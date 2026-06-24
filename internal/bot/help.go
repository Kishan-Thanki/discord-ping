package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) cmdHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := "!"
	if m.GuildID != "" {
		p, err := b.store.GetPrefix(context.Background(), m.GuildID)
		if err == nil {
			prefix = p
		}
	}

	embed := newEmbed("📚 Bot Commands", "Here is a list of all available commands.")

	categories := []struct {
		Name string
		Cmds []string
	}{
		{"Diagnostics", []string{"ping", "version", "uptime"}},
		{"Utility", []string{"help"}},
		{"Admin", []string{"setprefix [new]"}},
	}

	for _, cat := range categories {
		var formattedCmds []string
		for _, cmd := range cat.Cmds {
			formattedCmds = append(formattedCmds, "`"+prefix+cmd+"`")
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   cat.Name,
			Value:  strings.Join(formattedCmds, ", "),
			Inline: false,
		})
	}

	SendEmbed(s, m.ChannelID, embed)
}

func (b *Bot) cmdSetPrefix(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.GuildID == "" {
		embed := newEmbed("❌ Error", "This command can only be used in a server.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if !hasPermission(s, m.Author.ID, m.ChannelID, discordgo.PermissionAdministrator) {
		embed := newEmbed("🚫 Permission Denied", "You need Administrator permissions to change the prefix.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	if len(args) < 2 {
		embed := newEmbed("❌ Usage", "Usage: `!setprefix [new_prefix]`")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	newPrefix := args[1]
	err := b.store.SetPrefix(context.Background(), m.GuildID, newPrefix)
	if err != nil {
		embed := newEmbed("❌ Error", "Failed to update prefix in the database.")
		SendEmbed(s, m.ChannelID, embed)
		return
	}

	embed := newEmbed("✅ Prefix Updated", "The prefix for this server is now `"+newPrefix+"`")
	SendEmbed(s, m.ChannelID, embed)
}
