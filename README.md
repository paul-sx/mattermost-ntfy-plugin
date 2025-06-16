# mattermost-ntfy-plugin
This plugin allows users to setup ntfy notifications in response to messages in channels. Each channel can have its own topic to send messages to, but they all have to be to the same ntfy server.

## How to use

### 1. Build the Plugin

Run the following command to create a plugin package:

```sh
make dist
```

This will generate a distributable plugin file in the `dist/` directory.

### 2. Install the Plugin

1. Upload the generated plugin package to your Mattermost server.
2. Follow the Mattermost documentation to enable plugins and allow unsigned plugins if necessary.

### 3. Configure the Plugin

After uploading and activating the plugin:

- Open the plugin settings panel in Mattermost.
- Set the following options:
  - **Server**: The ntfy server URL.
  - **Topic**: The default topic for notifications.
  - **Username**: The username for ntfy authentication.
  - **Password**: The password for ntfy authentication.
  - **Enable**: Make sure to enable the notifications by toggling the last setting.

### 4. Using the `/ntfy` Slash Command

Once the plugin is enabled, you can use the `/ntfy` command in any channel:

- `/ntfy on` — Enable ntfy notifications for the current channel.
- `/ntfy off` — Disable ntfy notifications for the current channel.
- `/ntfy topic [topic]` — Change the broadcast topic for your user in the current channel.
- `/ntfy topic` — Reset the channel topic back to the default.

> **Tip:** Each channel can have its own topic, but all notifications go to the same ntfy server.

---

