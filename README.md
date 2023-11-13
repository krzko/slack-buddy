# slack-buddy

Slack Buddy is a convenient tool designed to help you manage your Slack status effortlessly. It operates as a system tray application, enabling you to quickly update your Slack status based on various pre-set or custom scenarios.

## Why

In the dynamic environment of Slack, keeping track of everyone's availability and status can be challenging, especially when mentioned inline in a conversation. Traditional status updates in Slack are often overlooked, particularly in fast-paced discussions. Slack Buddy addresses this gap effectively.

It updates not just your status emoji and text, but also your display name, embedding your current status directly within it. This enhancement makes your availability or activity instantly visible whenever you are mentioned or involved in a conversation, providing clearer context to your colleagues. This visibility is particularly crucial in remote or hybrid work settings where physical cues are absent, ensuring smoother, more informed interactions and collaborations on Slack.

### No Status

<img
  src="/assets/images/no-status.png"
  alt="No status set"
  title="No status set"
  style="display: inline-block; margin: 0 auto; max-width: 300px">

### With Status

<img
  src="/assets/images/with-status.png"
  alt="Slack Buddy status set"
  title="Slack Buddy status set"
  style="display: inline-block; margin: 0 auto; max-width: 300px">

## Features

* **Pre-defined Status Options**: Set your status to common activities like 'In a meeting', 'Commuting', 'Out sick', etc., with just a click.
* **Custom Status Support**: Add your personalized statuses with custom text and emojis.
* **Scheduled Status Updates**: Automate status updates for routine activities with start and end times.
* **Menu Customisation**: Organise your statuses into categories such as HR, Marketing, SRE, etc., for easy access.
* **Configuration Management**: Easily configure your Slack API token, User ID, and Display Name.

### navbar

<img
  src="/assets/images/navbar.png"
  alt="Slack Buddy Navigation Bar"
  title="Slack Buddy Navigation Bar"
  style="display: inline-block; margin: 0 auto; max-width: 300px">

## Install and Run

Install the package via on the of the below methods:

### brew

Install brew and then run:

```sh
brew install krzko/tap/slack-buddy
```

### Download Binary

Download the latest version from the [Releases](https://github.com/krzko/slack-buddy/releases) page.

### Run

Now that you have downloaded or installed the package, we need to run the binary. You can start `slack-buddy` via:

```sh
nohup slack-buddy > /tmp/slack-buddy.log &
```

A native installer is coming in the future.

## Usage

* **Initial Setup**: Enter your Slack API token, User ID, and Display Name through the settings menu.
* **Setting a Status**: Click on a predefined or custom status in the system tray menu to update your Slack status.
* **Adding Custom Statuses**: Use the settings menu to add or modify custom statuses, including their schedule if necessary.
* **Scheduling Statuses**: Define specific times for automatic status updates, perfect for regular activities like lunches or meetings.
* **Clearing Your Status**: Use the 'Clear Status' option to revert to your default state.

Slack Buddy is ideal for anyone who need to keep their team updated on their availability and activities without the hassle of manual updates. This tool is a must-have for enhancing your productivity and communication in a remote or hybrid work environment.

## Configuration File

Slack Buddy custom and scheduled statuses can be added via updating the`config.yaml` file in your home directory. An example config.yaml can be seen [here](./config_example.yaml). The file will be located in:

* **linux**: Soon
* **macOS**: `~/.slack-buddy/config.yaml`
* **Windows**: Soon

### `api_token`

You will need to retrieve a Personal Token from a Slack Workspace Administrator. You can get more details from the Slack [Access Tokens](https://api.slack.com/authentication/token-types) page.

### `user_id`

To retrieve youe `user_id` value, follow these steps:

* Open **Slack**.
* Open **Profile** window.
* Select the **three vertical dots** next to **View As**.
* Select **Copy memeber ID**.

### `disaply_name`

* This is the default disaply name value. All statuses will be set via **`display_name` + "is ..."**.
* This is the value that is set when you select **Clear Status**.

### Example Configuration

* The `title` value has to be unique.
* Custom statuses do not have `days`, `start_time` or `end_time` values set.
* Ensure the `status_emoji` value exists in your Slack workspace.

```yaml
api_token: xoxp-xxx-xxx-xxx-xxx
user_id: XXXXXXXXX
display_name: Kristof
custom_items:
  - title: "Zwift"
    tooltip: ""
    status_text: "on Zwift"
    status_emoji: ":bike:"
  - title: "School Run"
    tooltip: ""
    status_text: "on a school run"
    status_emoji: ":school_satchel:"
    days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    start_time: "08:45"
    end_time: "09:15"
  - title: "Daily Stand-up"
    tooltip: ""
    status_text: "in a meeting"
    status_emoji: ":spiral_calendar_pad:"
    days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    start_time: "10:00"
    end_time: "10:15"
  - title: "Lunch"
    tooltip: ""
    status_text: "lunching"
    status_emoji: ":pizza:"
    days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    start_time: "12:00"
    end_time: "13:00"
  - title: "School Run"
    tooltip: ""
    status_text: "on a school run"
    status_emoji: ":school_satchel:"
    days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    start_time: "15:15"
    end_time: "15:45"
```

## Rate Limits

Slack has some pretty limited rate limits, which can occur regularly.

> Profile update rate limits 
Update a user's profile, including custom status, sparingly. Special [rate limit](https://api.slack.com/docs/rate-limits) rules apply when updating profile data with `users.profile.set`. A token may update a single user's profile no more than **10** times per minute. And a single token may only set **30** user profiles per minute. Some burst behavior is allowed.

If your initial status update isn't set, wait a while and try again.
