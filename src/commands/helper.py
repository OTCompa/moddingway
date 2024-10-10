import discord
from settings import get_settings
from util import embed_field, create_embed

settings = get_settings()


def create_logging_embed(interaction: discord.Interaction):
    return create_embed(
        interaction.guild.get_channel(settings.logging_channel_id),
        user=interaction.user,
        timestamp=interaction.created_at,
        description=f"Used `{interaction.command.name}` command in {interaction.channel.mention}",
        fields=[embed_field("Action", f"/{interaction.command.name}")]
    )
