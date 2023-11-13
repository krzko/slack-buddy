package systray

import (
	"fmt"
	"log"

	"github.com/krzko/slack-buddy/pkg/config"
	"github.com/krzko/slack-buddy/pkg/slack"

	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/jasonlvhit/gocron"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Version = "0.0.1"

type SystrayHandler struct {
	slackClient *slack.SlackClient
	cfg         *config.Config
}

func OnReady() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	handler := &SystrayHandler{
		cfg: cfg,
	}
	handler.initializeSystray()

	// If configuration is first run, open settings dialog
	if handler.cfg.ApiToken == "your_default_api_token" || handler.cfg.UserId == "your_default_user_id" || handler.cfg.DisplayName == "your_default_display_name" {
		openSettingsDialog(cfg)
	}

	// After loading the config, check for duplicate titles
	checkForDuplicateTitles(cfg.CustomItems)
}

func (h *SystrayHandler) getSlackClient() (*slack.SlackClient, error) {
	if h.slackClient == nil {
		slackClient, err := slack.NewSlackClient(h.cfg.ApiToken, h.cfg.UserId, h.cfg.DisplayName)
		if err != nil {
			return nil, err
		}
		h.slackClient = slackClient
	}
	return h.slackClient, nil
}

func (h *SystrayHandler) initializeSystray() {
	// Start scheduling recurring statuses
	go h.scheduleRecurringStatuses()

	systray.SetIcon(Icon)

	inMeetingItem := systray.AddMenuItem("🗓 In a meeting", "")
	commutingItem := systray.AddMenuItem("🚗 Commuting", "")
	outSickItem := systray.AddMenuItem("🤕 Out sick", "")
	vactationingItem := systray.AddMenuItem("🏝 Vacationing", "")
	lunchingItem := systray.AddMenuItem("🍕 Lunching", "")
	onACallItem := systray.AddMenuItem("📲 On a call", "")

	systray.AddSeparator()

	workingRemitelyItem := systray.AddMenuItem("🏠 Working remotely", "")
	inOfficeItem := systray.AddMenuItem("🏢 Working in the office", "")

	systray.AddSeparator()

	hrItem := systray.AddMenuItem("🧑‍💼 Human Resources", "Manage HR-related activities")
	interviewingItem := hrItem.AddSubMenuItem("🤝 Interviewing", "")
	employeeOnboardingItem := hrItem.AddSubMenuItem("👨‍🏫 Employee Onboarding", "")

	marketingItem := systray.AddMenuItem("📈 Marketing", "Manage marketing activities")
	campaignPlanningItem := marketingItem.AddSubMenuItem("📊 Campaign Planning", "")
	socialMediaUpdateItem := marketingItem.AddSubMenuItem("📱 Social Media Update", "")

	personalLifeItem := systray.AddMenuItem("🏠 Personal Life", "Statuses for personal activities")
	doctorAppointmentItem := personalLifeItem.AddSubMenuItem("👨‍⚕️ Doctor Appointment", "")
	gamingItem := personalLifeItem.AddSubMenuItem("🎮 Gaming", "Taking a break to play games")
	readingItem := personalLifeItem.AddSubMenuItem("📖 Reading", "Taking some time to read")
	runningErrandsItem := personalLifeItem.AddSubMenuItem("🛒 Running Errands", "")
	schoolRunItem := personalLifeItem.AddSubMenuItem("🎒 School Run", "Doing the shcool run")
	travelingItem := personalLifeItem.AddSubMenuItem("✈️ Traveling", "")

	sreItem := systray.AddMenuItem("👩‍🚒 Site Reliability Engineering", "Site Reliability Engineering activities")
	incidentCommanderItem := sreItem.AddSubMenuItem("👩‍🚒 Incident Commander (IC)", "")
	communicationsLeadItem := sreItem.AddSubMenuItem("📡 Communications Lead (CL)", "")
	operationsLeadItem := sreItem.AddSubMenuItem("🚀 Operations Lead (OL)", "")

	softwareDevItem := systray.AddMenuItem("💻 Software Development", "Software development activities")
	onCallDevItem := softwareDevItem.AddSubMenuItem("👨‍💻 On Call Developer", "")
	codeReviewItem := softwareDevItem.AddSubMenuItem("🔍 Code Review", "")

	// Initialize a map that uses pointers to systray.MenuItem as keys and config.CustomItem as values
	customMenuItems := make(map[*systray.MenuItem]config.CustomItem)

	// Add custom items menu
	if len(h.cfg.CustomItems) > 0 {
		systray.AddSeparator()
		customItem := systray.AddMenuItem("🆕 Custom", "Your custom statuses")
		scheduledItem := systray.AddMenuItem("🕒 Scheduled", "Your scheduled statuses")

		// Populate the map with menu items
		for _, item := range h.cfg.CustomItems {
			var menuItemToAdd *systray.MenuItem
			if item.StartTime != "" && item.EndTime != "" {
				// Format the title to include the start and end times
				titleWithTime := fmt.Sprintf("%s at %s - %s", item.Title, item.StartTime, item.EndTime)
				menuItemToAdd = scheduledItem.AddSubMenuItem(titleWithTime, item.Tooltip)
				menuItemToAdd.Disable()
			} else {
				menuItemToAdd = customItem.AddSubMenuItem(item.Title, item.Tooltip)
			}
			customMenuItems[menuItemToAdd] = item
		}
	}

	systray.AddSeparator()

	clearStatusItem := systray.AddMenuItem("🧹 Clear Status", "")

	systray.AddSeparator()

	settingsItem := systray.AddMenuItem("⚙️ Settings", "")
	quitItem := systray.AddMenuItem("❌ Quit", "")

	go func() {
		for {
			select {
			case <-inMeetingItem.ClickedCh:
				err := h.updateStatus("in a meeting", ":spiral_calendar_pad:")
				if err != nil {
					// Handle error
				}

			case <-commutingItem.ClickedCh:
				err := h.updateStatus("commuting", ":car:")
				if err != nil {
					// Handle error
				}
			case <-outSickItem.ClickedCh:
				err := h.updateStatus("out sick", ":face_with_thermometer:")
				if err != nil {
					// Handle error
				}
			case <-vactationingItem.ClickedCh:
				err := h.updateStatus("vacationing", ":palm_tree:")
				if err != nil {
					// Handle error
				}
			case <-lunchingItem.ClickedCh:
				err := h.updateStatus("lunching", ":pizza:")
				if err != nil {
					// Handle error
				}
			case <-onACallItem.ClickedCh:
				err := h.updateStatus("on a call", ":calling:")
				if err != nil {
					// Handle error
				}

			case <-workingRemitelyItem.ClickedCh:
				err := h.updateStatus("working remotely", ":office:")
				if err != nil {
					// Handle error
				}
			case <-inOfficeItem.ClickedCh:
				err := h.updateStatus("working in the office", ":house:")
				if err != nil {
					// Handle error
				}

			case <-incidentCommanderItem.ClickedCh:
				err := h.updateStatus("Incident Commander (IC)", ":firefighter:")
				if err != nil {
					// Handle error
				}
			case <-communicationsLeadItem.ClickedCh:
				err := h.updateStatus("Communications Lead (CL)", ":mega:")
				if err != nil {
					// Handle error
				}
			case <-operationsLeadItem.ClickedCh:
				err := h.updateStatus("Operations Lead (CL)", ":technologist:")
				if err != nil {
					// Handle error
				}

			case <-onCallDevItem.ClickedCh:
				err := h.updateStatus("the on-call developer", ":male-technologist:")
				if err != nil {
					// Handle error
				}
			case <-codeReviewItem.ClickedCh:
				err := h.updateStatus("performing code reviews", ":magnifying_glass_tilted_left:")
				if err != nil {
					// Handle error
				}
			case <-campaignPlanningItem.ClickedCh:
				err := h.updateStatus("campaign planning", ":bar_chart:")
				if err != nil {
					// Handle error
				}
			case <-socialMediaUpdateItem.ClickedCh:
				err := h.updateStatus("doing social media updates", ":iphone:")
				if err != nil {
					// Handle error
				}
			case <-interviewingItem.ClickedCh:
				err := h.updateStatus("interviewing", ":handshake:")
				if err != nil {
					// Handle error
				}
			case <-employeeOnboardingItem.ClickedCh:
				err := h.updateStatus("employee onboarding", ":teacher:")
				if err != nil {
					// Handle error
				}
			case <-doctorAppointmentItem.ClickedCh:
				err := h.updateStatus("at a doctor's appointment", ":male-doctor:")
				if err != nil {
					// Handle error
				}
			case <-travelingItem.ClickedCh:
				err := h.updateStatus("traveling", ":airplane:")
				if err != nil {
					// Handle error
				}
			case <-runningErrandsItem.ClickedCh:
				err := h.updateStatus("running errands", ":shopping_cart:")
				if err != nil {
					// Handle error
				}
			case <-gamingItem.ClickedCh:
				err := h.updateStatus("gaming", ":video_game:")
				if err != nil {
					// Handle error
				}
			case <-readingItem.ClickedCh:
				err := h.updateStatus("reading", ":book:")
				if err != nil {
					// Handle error
				}
			case <-schoolRunItem.ClickedCh:
				err := h.updateStatus("on a school run", ":school_satchel:")
				if err != nil {
					// Handle error
				}

			case <-clearStatusItem.ClickedCh:
				err := h.unsetStatus()
				if err != nil {
					// Handle error
				}

			case <-settingsItem.ClickedCh:
				// fmt.Println("Settings item clicked")
				openSettingsDialog(h.cfg)

			case <-quitItem.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// New goroutine for handling custom menu items and separating scheduled items
	if len(h.cfg.CustomItems) > 0 {
		go func() {
			for {
				for menuItem, customItem := range customMenuItems {
					<-menuItem.ClickedCh
					err := h.updateStatus(customItem.StatusText, customItem.StatusEmoji)
					if err != nil {
						// Handle error
					}
				}
			}
		}()
	}
}

func OnExit() {}

func openSettingsDialog(cfg *config.Config) {
	apiToken, _, err := dlgs.Entry("API Token", "Enter your Slack API Token:", cfg.ApiToken)
	if err != nil {
		fmt.Printf("Error displaying API Token dialog: %v\n", err)
		return
	}
	cfg.ApiToken = apiToken

	userId, _, err := dlgs.Entry("User ID", "Enter your Slack User ID:", cfg.UserId)
	if err != nil {
		fmt.Printf("Error displaying User ID dialog: %v\n", err)
		return
	}
	cfg.UserId = userId

	displayName, _, err := dlgs.Entry("Display Name", "Enter your Slack Display Name:", cfg.DisplayName)
	if err != nil {
		fmt.Printf("Error displaying Display Name dialog: %v\n", err)
		return
	}
	cfg.DisplayName = displayName

	// Save configuration
	err = cfg.SaveConfig()
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		return
	}
}

func createJob(day, startTime, endTime string, updateFunc, unsetFunc interface{}) {
	job := gocron.Every(1)
	switch day {
	case "Monday":
		job.Monday().At(startTime).Do(updateFunc)
	case "Tuesday":
		job.Tuesday().At(startTime).Do(updateFunc)
	case "Wednesday":
		job.Wednesday().At(startTime).Do(updateFunc)
	case "Thursday":
		job.Thursday().At(startTime).Do(updateFunc)
	case "Friday":
		job.Friday().At(startTime).Do(updateFunc)
	case "Saturday":
		job.Saturday().At(startTime).Do(updateFunc)
	case "Sunday":
		job.Sunday().At(startTime).Do(updateFunc)

	}
	gocron.Every(1).Day().At(endTime).Do(unsetFunc)
}

func (h *SystrayHandler) scheduleRecurringStatuses() {
	for _, item := range h.cfg.CustomItems {
		for _, day := range item.Days {
			createJob(day, item.StartTime, item.EndTime, func() { h.updateStatus(item.StatusText, item.StatusEmoji) }, h.unsetStatus)
		}
	}
	<-gocron.Start()
}

func (h *SystrayHandler) updateStatus(statusText, statusEmoji string) error {
	caser := cases.Title(language.English)
	capitalisedStatusText := caser.String(statusText)

	fmt.Printf("SystrayHandler updateStatus called with statusText=%s, statusEmoji=%s\n", statusText, statusEmoji)
	slackClient, err := h.getSlackClient()
	if err != nil {
		fmt.Printf("Error getting Slack client: %v\n", err)
		return err
	}

	stdLibClient := slack.NewStdLibClient(h.cfg.ApiToken)
	newDisplayName := fmt.Sprintf("%s is %s", h.cfg.DisplayName, statusText)
	err = stdLibClient.UpdateDisplayName(h.cfg.UserId, newDisplayName)
	if err != nil {
		fmt.Printf("Error updating display name: %v\n", err)
		return err
	}

	err = slackClient.UpdateStatus(capitalisedStatusText, statusEmoji)
	if err != nil {
		fmt.Printf("Error updating status: %v\n", err)
		return err
	}
	return nil
}

func (h *SystrayHandler) unsetStatus() error {
	stdLibClient := slack.NewStdLibClient(h.cfg.ApiToken)
	err := stdLibClient.UpdateDisplayName(h.cfg.UserId, h.cfg.DisplayName)
	if err != nil {
		fmt.Printf("Error updating display name: %v\n", err)
		return err
	}

	slackClient, err := h.getSlackClient()
	if err != nil {
		return err
	}
	return slackClient.UnsetStatus()
}

func findCustomItemByStatusText(items []config.CustomItem, statusText string) *config.CustomItem {
	for _, item := range items {
		if item.StatusText == statusText {
			return &item
		}
	}
	return nil
}

func findCustomItemByTitle(items []config.CustomItem, title string) *config.CustomItem {
	for _, item := range items {
		if item.Title == title {
			return &item
		}
	}
	return nil
}

func checkForDuplicateTitles(items []config.CustomItem) {
	titles := make(map[string]bool)
	for _, item := range items {
		if _, exists := titles[item.Title]; exists {
			log.Fatalf("duplicate title found: %s", item.Title)
		}
		titles[item.Title] = true
	}
}
