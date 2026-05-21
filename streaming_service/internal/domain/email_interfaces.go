package domain

type EmailSender interface {
	SendPlaylistSharedNotification(toEmail, playlistName, senderName string) error
	SendTrendingSongsNotification(toEmail string) error
	SendMilestoneNotification(toEmail string, songsCount int) error
}
