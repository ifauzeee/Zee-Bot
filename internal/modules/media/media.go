package media

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zee-ubot/internal/handlers"

	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func Register(m *handlers.Manager) {
	m.Register("save", "Save view once or timer media", saveHandler)
	m.Register("vo", "Alias for save (Get View Once)", saveHandler)
}

func saveHandler(c *handlers.Context) error {
	if ch, ok := c.Peer.(*tg.InputPeerChannel); ok {
		_, _ = c.Raw.ChannelsDeleteMessages(c.Ctx, &tg.ChannelsDeleteMessagesRequest{
			Channel: &tg.InputChannel{ChannelID: ch.ChannelID, AccessHash: ch.AccessHash},
			ID:      []int{c.Msg.ID},
		})
	} else {
		_, _ = c.Raw.MessagesDeleteMessages(c.Ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     []int{c.Msg.ID},
		})
	}

	reply, err := getRepliedMessage(c)
	if err != nil {
		c.Logger.Error("save: getRepliedMessage failed", zap.Error(err))
		return nil
	}

	if reply.Media == nil {
		return nil
	}

	startTime := time.Now()

	d := downloader.NewDownloader()
	path, isVideo, err := downloadMedia(c.Ctx, c.Raw, d, reply.Media)
	if err != nil {
		c.Logger.Error("Download failed", zap.Error(err))
		return nil
	}
	defer os.Remove(path)

	duration := time.Since(startTime).Round(time.Millisecond)

	u := uploader.NewUploader(c.Raw)
	upload, err := u.FromPath(c.Ctx, path)
	if err != nil {
		c.Logger.Error("Upload failed", zap.Error(err))
		return nil
	}

	target := c.Sender.Self()
	var sendErr error

	caption := fmt.Sprintf("🔓 <b>Media Saved</b>\nTook: %s", duration)

	ext := strings.ToLower(filepath.Ext(path))
	if isVideo {
		_, sendErr = target.Video(c.Ctx, upload, html.String(nil, caption))
	} else if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		_, sendErr = target.UploadedPhoto(c.Ctx, upload, html.String(nil, caption))
	} else {
		_, sendErr = target.File(c.Ctx, upload, html.String(nil, caption))
	}

	if sendErr != nil {
		c.Logger.Error("Send to Saved Messages failed", zap.Error(sendErr))
	}

	return nil
}

func getRepliedMessage(c *handlers.Context) (*tg.Message, error) {
	if c.Msg.ReplyTo == nil {
		return nil, fmt.Errorf("reply to a message")
	}

	replyHeader, ok := c.Msg.ReplyTo.(*tg.MessageReplyHeader)
	if !ok {
		return nil, fmt.Errorf("invalid reply")
	}

	id := replyHeader.ReplyToMsgID

	var result tg.MessagesMessagesClass
	var err error

	if inputPeerChannel, ok := c.Peer.(*tg.InputPeerChannel); ok {
		inputChannel := &tg.InputChannel{
			ChannelID:  inputPeerChannel.ChannelID,
			AccessHash: inputPeerChannel.AccessHash,
		}
		result, err = c.Raw.ChannelsGetMessages(c.Ctx, &tg.ChannelsGetMessagesRequest{
			Channel: inputChannel,
			ID:      []tg.InputMessageClass{&tg.InputMessageID{ID: id}},
		})
	} else {
		result, err = c.Raw.MessagesGetMessages(c.Ctx, []tg.InputMessageClass{&tg.InputMessageID{ID: id}})
	}

	if err != nil {
		return nil, err
	}

	var messages []tg.MessageClass
	switch m := result.(type) {
	case *tg.MessagesMessages:
		messages = m.Messages
	case *tg.MessagesMessagesSlice:
		messages = m.Messages
	case *tg.MessagesChannelMessages:
		messages = m.Messages
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("message not found")
	}

	msg, ok := messages[0].(*tg.Message)
	if !ok {
		return nil, fmt.Errorf("target is empty or service message")
	}

	return msg, nil
}

func downloadMedia(ctx context.Context, api *tg.Client, d *downloader.Downloader, media tg.MessageMediaClass) (string, bool, error) {
	tmpDir := "tmp"
	_ = os.MkdirAll(tmpDir, 0755)

	isVideo := false
	var outPath string
	var location tg.InputFileLocationClass

	switch m := media.(type) {
	case *tg.MessageMediaPhoto:
		p, ok := m.Photo.(*tg.Photo)
		if !ok {
			return "", false, fmt.Errorf("photo content missing")
		}

		var bestSize string
		var maxW int
		for _, s := range p.Sizes {
			if ps, ok := s.(*tg.PhotoSize); ok {
				if ps.W > maxW {
					maxW = ps.W
					bestSize = ps.Type
				}
			} else if ps, ok := s.(*tg.PhotoSizeProgressive); ok {
				if ps.W > maxW {
					maxW = ps.W
					bestSize = ps.Type
				}
			}
		}

		if bestSize == "" && len(p.Sizes) > 0 {
			switch s := p.Sizes[len(p.Sizes)-1].(type) {
			case *tg.PhotoSize:
				bestSize = s.Type
			case *tg.PhotoSizeProgressive:
				bestSize = s.Type
			}
		}

		location = &tg.InputPhotoFileLocation{
			ID:            p.ID,
			AccessHash:    p.AccessHash,
			FileReference: p.FileReference,
			ThumbSize:     bestSize,
		}
		outPath = filepath.Join(tmpDir, fmt.Sprintf("photo_%d.jpg", p.ID))

	case *tg.MessageMediaDocument:
		doc, ok := m.Document.(*tg.Document)
		if !ok {
			return "", false, fmt.Errorf("document content missing")
		}

		location = &tg.InputDocumentFileLocation{
			ID:            doc.ID,
			AccessHash:    doc.AccessHash,
			FileReference: doc.FileReference,
			ThumbSize:     "",
		}

		fname := fmt.Sprintf("doc_%d", doc.ID)
		for _, attr := range doc.Attributes {
			if fn, ok := attr.(*tg.DocumentAttributeFilename); ok {
				fname = fn.FileName
			}
			if _, ok := attr.(*tg.DocumentAttributeVideo); ok {
				isVideo = true
			}
		}

		if !strings.Contains(fname, ".") {
			if strings.Contains(doc.MimeType, "video") {
				fname += ".mp4"
				isVideo = true
			} else if strings.Contains(doc.MimeType, "image") {
				fname += ".jpg"
			}
		}

		outPath = filepath.Join(tmpDir, fname)

	default:
		return "", false, fmt.Errorf("unsupported media type")
	}

	_, err := d.Download(api, location).ToPath(ctx, outPath)
	if err != nil {
		return "", false, err
	}

	return outPath, isVideo, nil
}
