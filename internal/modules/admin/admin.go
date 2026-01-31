package admin

import (
	"zee-ubot/internal/handlers"

	"github.com/gotd/td/tg"
)

func Register(m *handlers.Manager) {
	m.Register("purge", "Delete messages", purgeHandler)
	m.Register("kick", "Kick user from group", kickHandler)
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

	switch m := msgs.(type) {
	case *tg.MessagesChannelMessages:
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
		return c.Edit("Failed to kick user: " + err.Error())
	}

	return c.Edit("User kicked!")
}
