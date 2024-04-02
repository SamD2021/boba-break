# Boba Break

Boba Break is a terminal-based application designed to help you manage your work-break cycles efficiently. It provides features such as break management, note-taking, and a main menu interface to navigate between different functionalities.

## Installation

To install Boba Break, follow these steps:

1. Clone the repository to your local machine:

```
git clone https://github.com/SamD2021/boba-break.git
```

2. Navigate to the project directory:

```
cd boba-break
```

3. Build the application:

```
go build
```

4. Run the application:

```
./boba-break
```

## Features

### Break Manager

The Break Manager module allows you to set work and break durations. It displays a timer indicating the time remaining for your work session. When the work session ends, it prompts you to take a break, and vice versa. You can control the timer using keyboard shortcuts.

### Main Menu

The Main Menu module provides a menu interface to access different features of the application. It currently supports navigation to the Break Manager and Notes modules.

### Notes

The Notes module allows you to jot down your thoughts or important information during work sessions. It provides a simple text editor interface with basic editing functionalities.

### Break Log

The Break Log module is a work-in-progress feature intended to log your break activities and durations. It currently supports adding log entries to a JSON file.

## Usage

Upon launching the application, you will be presented with the main menu. From there, you can navigate to the Break Manager to start your work-break cycles or to the Notes module to take notes. Use the provided keyboard shortcuts to control the timer and navigate through the application.

## Roadmap

### Version 1.0
- [x] Implement basic Break Manager UI with timer functionality.
- [x] Add support for starting, pausing, and resetting timers.
- [x] Integrate main menu UI for navigation between features.

### Version 1.1
- [ ] Implement customizable work and break durations.
- [ ] Add sound notifications for timer events.
- [ ] Integrate visual indicators for timer progress.

### Version 1.2
- [ ] Implement Notes UI for taking and saving notes.
- [x] Add support for basic text editing functionalities.
- [ ] Integrate saving and loading notes from disk.

### Version 1.3
- [x] Research and plan daemon implementation for background timer functionality.
- [ ] Define communication protocols between UI and daemon components.
- [ ] Implement daemon functionality to handle timer logic in the background.

### Version 1.4
- [ ] Write comprehensive documentation for installation, usage, and contribution guidelines.
- [ ] Prepare project for deployment by packaging necessary files and dependencies.
- [ ] Create release tags and versioning for stable releases.

### Version 1.5 (Future Enhancements)
- [ ] Explore additional features such as statistics tracking and integration with external tools.
- [ ] Consider platform-specific optimizations or adaptations for different operating systems.
- [ ] Experiment with UI/UX improvements based on user feedback and industry trends.

## Contributing

Contributions to Boba Break are welcome! If you find any bugs or have suggestions for new features, feel free to open an issue or submit a pull request on GitHub.

## License

This project is licensed under the GPL 3.0 License - see the [LICENSE](LICENSE) file for details.

---
