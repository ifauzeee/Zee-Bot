package admin

import (
	"fmt"
	"strings"
	"time"
	"zee-ubot/internal/handlers"

	"github.com/gotd/td/tg"
)

func init() {
	handlers.RegisterModule(Register)
}

func Register(m *handlers.Manager) {
	m.Register("purge", "Delete messages", purgeHandler)
	m.Register("kick", "Kick user from group", kickHandler)
	m.Register("zombies", "Clean deleted accounts from group", zombiesHandler)
	m.Register("info", "Get detailed info of a user or chat", infoHandler)
}

func purgeHandler(c *handlers.Context) error {
	replyHeader, ok := c.Msg.ReplyTo.(*tg.MessageReplyHeader)
	if !ok || replyHeader == nil {
		return c.Edit("Please reply to a message to start purging.")
	}

	replyID := replyHeader.ReplyToMsgID
	currentID := c.Msg.ID

	var ids []int
	for i := replyID; i <= currentID; i++ {
		ids = append(ids, i)
	}

	_, err := c.Raw.MessagesDeleteMessages(c.Ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     ids,
		Revoke: true,
	})
	return err
}

func kickHandler(c *handlers.Context) error {
	replyHeader, ok := c.Msg.ReplyTo.(*tg.MessageReplyHeader)
	if !ok || replyHeader == nil {
		return c.Edit("Please reply to a user to kick.")
	}

	inputChannel, ok := c.Peer.(*tg.InputPeerChannel)
	if !ok {
		return c.Edit("This command can only be used in channels/supergroups.")
	}

	msgs, err := c.Raw.ChannelsGetMessages(c.Ctx, &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  inputChannel.ChannelID,
			AccessHash: inputChannel.AccessHash,
		},
		ID: []tg.InputMessageClass{&tg.InputMessageID{ID: replyHeader.ReplyToMsgID}},
	})

	if err != nil {
		return c.Edit("Failed to fetch replied message: " + err.Error())
	}

	var targetUserID int64
	var targetAccessHash int64
	found := false

	if m, ok := msgs.(*tg.MessagesChannelMessages); ok {
		for _, msg := range m.Messages {
			if msgVal, ok := msg.(*tg.Message); ok {
				if peerUser, ok := msgVal.FromID.(*tg.PeerUser); ok {
					targetUserID = peerUser.UserID
					for _, u := range m.Users {
						if user, ok := u.(*tg.User); ok && user.ID == targetUserID {
							targetAccessHash = user.AccessHash
							found = true
							break
						}
					}
				}
			}
		}
	}

	if !found {
		return c.Edit("Could not find user to kick.")
	}

	_, err = c.Raw.ChannelsEditBanned(c.Ctx, &tg.ChannelsEditBannedRequest{
		Channel: &tg.InputChannel{
			ChannelID:  inputChannel.ChannelID,
			AccessHash: inputChannel.AccessHash,
		},
		Participant: &tg.InputPeerUser{
			UserID:     targetUserID,
			AccessHash: targetAccessHash,
		},
		BannedRights: tg.ChatBannedRights{
			ViewMessages: true,
		},
	})

	if err != nil {
		style := c.NewStyle("Kick Failed", "❌")
		style.AddRow("Error", err.Error())
		return c.Edit(style.Build())
	}

	style := c.NewStyle("Kick Success", "✅")
	style.AddRow("User Status", "Successfully kicked from group")
	return c.Edit(style.Build())
}

func zombiesHandler(c *handlers.Context) error {
	inputChannel, ok := c.Peer.(*tg.InputPeerChannel)
	if !ok {
		return c.EditStatus("❌", "This command only works in supergroups.")
	}

	_ = c.EditStatus("🔍", "Searching for zombies...")

	var zombies []*tg.User
	offset := 0
	limit := 100

	for {
		res, err := c.Raw.ChannelsGetParticipants(c.Ctx, &tg.ChannelsGetParticipantsRequest{
			Channel: &tg.InputChannel{
				ChannelID:  inputChannel.ChannelID,
				AccessHash: inputChannel.AccessHash,
			},
			Filter: &tg.ChannelParticipantsRecent{},
			Offset: offset,
			Limit:  limit,
		})

		if err != nil {
			return c.Edit("❌ Failed to get participants: " + err.Error())
		}

		p, ok := res.(*tg.ChannelsChannelParticipants)
		if !ok {
			break
		}

		for _, u := range p.Users {
			if user, ok := u.(*tg.User); ok && user.Deleted {
				zombies = append(zombies, user)
			}
		}

		if len(p.Participants) < limit {
			break
		}
		offset += limit
	}

	if len(zombies) == 0 {
		style := c.NewStyle("Zombie Cleanup", "🧹")
		style.AddRow("Status", "No zombies found in this group")
		return c.Edit(style.Build())
	}

	_ = c.EditStatus("🧟", "Found "+fmt.Sprintf("%d", len(zombies))+" zombies. Cleaning up...")

	count := 0
	for _, z := range zombies {
		_, err := c.Raw.ChannelsEditBanned(c.Ctx, &tg.ChannelsEditBannedRequest{
			Channel: &tg.InputChannel{
				ChannelID:  inputChannel.ChannelID,
				AccessHash: inputChannel.AccessHash,
			},
			Participant: &tg.InputPeerUser{
				UserID:     z.ID,
				AccessHash: z.AccessHash,
			},
			BannedRights: tg.ChatBannedRights{
				ViewMessages: true,
			},
		})
		if err == nil {
			count++
		}
		time.Sleep(200 * time.Millisecond)
	}

	style := c.NewStyle("Zombie Cleanup", "🧹")
	style.AddRow("Zombies Removed", count)
	style.AddRow("Status", "Group is now clean")

	return c.Edit(style.Build())
}

func infoHandler(c *handlers.Context) error {
	var targetID int64
	var accessHash int64
	var username string
	var peerType string
	var title string

	if c.Msg.ReplyTo != nil {
		replyID := c.Msg.ReplyTo.(*tg.MessageReplyHeader).ReplyToMsgID
		var result tg.MessagesMessagesClass
		var err error

		if ch, ok := c.Peer.(*tg.InputPeerChannel); ok {
			result, err = c.Raw.ChannelsGetMessages(c.Ctx, &tg.ChannelsGetMessagesRequest{
				Channel: &tg.InputChannel{ChannelID: ch.ChannelID, AccessHash: ch.AccessHash},
				ID:      []tg.InputMessageClass{&tg.InputMessageID{ID: replyID}},
			})
		} else {
			result, err = c.Raw.MessagesGetMessages(c.Ctx, []tg.InputMessageClass{&tg.InputMessageID{ID: replyID}})
		}

		if err == nil {
			var users []tg.UserClass
			switch m := result.(type) {
			case *tg.MessagesMessages:
				users = m.Users
			case *tg.MessagesMessagesSlice:
				users = m.Users
			case *tg.MessagesChannelMessages:
				users = m.Users
			}

			if len(users) > 0 {
				if user, ok := users[0].(*tg.User); ok {
					targetID = user.ID
					accessHash = user.AccessHash
					username = user.Username
					title = strings.TrimSpace(user.FirstName + " " + user.LastName)
					peerType = "User"
				}
			}
		}
	} else {
		switch p := c.Peer.(type) {
		case *tg.InputPeerChannel:
			res, _ := c.Raw.ChannelsGetChannels(c.Ctx, []tg.InputChannelClass{&tg.InputChannel{
				ChannelID:  p.ChannelID,
				AccessHash: p.AccessHash,
			}})
			if chats := res.GetChats(); len(chats) > 0 {
				if ch, ok := chats[0].(*tg.Channel); ok {
					targetID = ch.ID
					accessHash = ch.AccessHash
					username = ch.Username
					peerType = "Channel/Supergroup"
					title = ch.Title
				}
			}
		case *tg.InputPeerUser:
			targetID = p.UserID
			accessHash = p.AccessHash
			peerType = "User (Direct)"
		case *tg.InputPeerChat:
			targetID = p.ChatID
			peerType = "Chat (Group)"
		}
	}

	if targetID == 0 {
		return c.Edit("❌ Could not get info for this entity.")
	}

	if title == "" {
		title = "Unknown"
	}
	if username == "" {
		username = "None"
	}

	style := c.NewStyle("Entity Information", "📊")
	style.AddRowWithIcon("🏷", "Title/Name", title)
	style.AddRowWithIcon("🆔", "ID", targetID)
	style.AddRowWithIcon("🔑", "AccessHash", accessHash)
	style.AddRowWithIcon("👤", "Username", "@"+username)
	style.AddRowWithIcon("🛠", "Type", peerType)

	return c.Edit(style.Build())
}
