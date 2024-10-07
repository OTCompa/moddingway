import discord
from discord.ext.commands import Bot
from settings import get_settings
from services.exile_service import exile_user, unexile_user
from util import user_has_role, calculate_time_delta
from enums import Role
from typing import Optional
from .helper import create_logging_embed

settings = get_settings()


async def is_user_moderator(interaction: discord.Interaction):
    return user_has_role(interaction.user, Role.ADMINISTRATION) or user_has_role(
        interaction.user, Role.MANAGEMENT
    )


def create_exile_commands(bot: Bot) -> None:
    @bot.tree.command()
    @discord.app_commands.check(is_user_moderator)
    @discord.app_commands.describe(user="User being exiled")
    async def unexile(interaction: discord.Interaction, user: discord.Member):
        """Unexile the specified user."""

        async with create_logging_embed(interaction) as logging_embed:
            await unexile_user(logging_embed, user)

            await interaction.response.send_message(
                f"Successfully unexiled {user.mention}", ephemeral=True
            )

    @bot.tree.command()
    @discord.app_commands.check(is_user_moderator)
    @discord.app_commands.describe(
        user="User being exiled",
        duration="The duration of the exile. TBA format",
        reason="Reason for exile",
    )
    async def exile(
        interaction: discord.Interaction,
        user: discord.Member,
        duration: Optional[str],
        reason: str,
    ):
        """Exile the specified user."""

        async with create_logging_embed(interaction) as logging_embed:
            exile_duration = calculate_time_delta(duration)
            if duration and not exile_duration:
                await interaction.response.send_message(
                    "Invalid exile duration given, duration should be in the form of [1 or 2 digits][s, d, m, h]",
                    ephemeral=True,
                )
                return
            await exile_user(logging_embed, user, exile_duration, reason)

            await interaction.response.send_message(
                f"Successfully exiled {user.mention}", ephemeral=True
            )
