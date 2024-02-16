package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v59/github"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type Version struct {
	Version              string `json:"version"`
	Date                 string `json:"date"`
	LocalizedDescription string `json:"localizedDescription"`
	DownloadURL          string `json:"downloadURL,omitempty"`
	Size                 int    `json:"size,omitempty"`
}

type Permission struct {
	PermissionType   string `json:"type"`
	UsageDescription string `json:"usageDescription"`
}

type Entitlement struct {
	Name string `json:"name"`
}

type App struct {
	Name                 string       `json:"name"`
	BundleIdentifier     string       `json:"bundleIdentifier"`
	DeveloperName        string       `json:"developerName"`
	Versions             []Version    `json:"versions"`
	Version              string       `json:"version,omitempty"`
	VersionDate          string       `json:"versionDate,omitempty"`
	VersionDescription   string       `json:"versionDescription,omitempty"`
	DownloadURL          string       `json:"downloadURL,omitempty"`
	LocalizedDescription string       `json:"localizedDescription"`
	IconURL              string       `json:"iconURL"`
	TintColor            string       `json:"tintColor"`
	Size                 int          `json:"size,omitempty"`
	Screenshots          []string     `json:"screenshots,omitempty"`
	Permissions          []Permission `json:"permissions"`
	AppPermissions       struct {
		Entitlements []Entitlement `json:"entitlements"`
	} `json:"appPermissions"`
}

type Source struct {
	Name        string `json:"name"`
	Identifier  string `json:"identifier"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Apps        []App  `json:"apps"`
}

var AidokuApp = App{
	Name:                 "Aidoku",
	BundleIdentifier:     "xyz.skitty.Aidoku",
	DeveloperName:        "Skitty",
	Versions:             []Version{},
	LocalizedDescription: "Free and open source manga reader for iOS and iPadOS",
	IconURL:              "https://avatars.githubusercontent.com/u/97767528",
	TintColor:            "ff2f52",
	Permissions: []Permission{
		{
			PermissionType:   "background-fetch",
			UsageDescription: "Aidoku periodically updates sources in the background.",
		},
		{
			PermissionType:   "background-audio",
			UsageDescription: "Allows Aidoku to run longer than 30 seconds when refreshing apps in background.",
		},
	},
	AppPermissions: struct {
		Entitlements []Entitlement `json:"entitlements"`
	}{
		Entitlements: []Entitlement{
			{
				Name: "get-task-allow",
			},
			{
				Name: "com.apple.security.application-groups",
			},
			{
				Name: "aps-environment",
			},
			{
				Name: "com.apple.developer.siri",
			},
		},
	},
	Screenshots: []string{
		"https://raw.githubusercontent.com/Caslus/AidokuSource/dev/media/screenshot1.jpeg",
		"https://raw.githubusercontent.com/Caslus/AidokuSource/dev/media/screenshot2.jpeg",
		"https://raw.githubusercontent.com/Caslus/AidokuSource/dev/media/screenshot3.jpeg",
	},
}

func updateApp() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(nil).WithAuthToken(githubToken)
	releases, res, err := client.Repositories.ListReleases(context.Background(), "Aidoku", "Aidoku", nil)
	if err != nil {
		log.Fatal("Error getting releases: ", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal("Error getting releases: ", res.Status)
	}

	AidokuSourceFile, err := os.ReadFile("source.json")
	if err != nil {
		log.Fatal("Error opening source file: ", err)
	}
	AidokuSourceContent := Source{}
	err = json.Unmarshal(AidokuSourceFile, &AidokuSourceContent)
	if err != nil {
		log.Fatal("Error unmarshalling source file: ", err)
	}
	if booting {
		booting = false
		AidokuSource = AidokuSourceContent
	}
	if AidokuSourceContent.Apps[0].Versions[0].Version == *releases[0].Name {
		fmt.Println("SOURCE ALREADY UP TO DATE! LATEST VERSION: ", *releases[0].Name)
		return
	}
	fmt.Println("UPDATING SOURCE...")
	var versions []Version
	for _, release := range releases {
		versions = append(versions, Version{
			Version:              *release.Name,
			Date:                 release.PublishedAt.Format(time.RFC3339),
			LocalizedDescription: *release.Body,
			DownloadURL:          *release.Assets[0].BrowserDownloadURL,
			Size:                 *release.Assets[0].Size,
		})
	}
	AidokuApp.Versions = versions
	AidokuApp.Version = versions[0].Version
	AidokuApp.VersionDate = versions[0].Date
	AidokuApp.VersionDescription = versions[0].LocalizedDescription
	AidokuApp.DownloadURL = versions[0].DownloadURL
	AidokuApp.Size = versions[0].Size
	AidokuSource.Apps[0] = AidokuApp
	newJson, err := json.MarshalIndent(AidokuSource, "", "  ")
	err = os.WriteFile("source.json", newJson, 0644)
	if err != nil {
		log.Fatal("Error writing source file: ", err)
	}
	fmt.Println("SOURCE UPDATED! LATEST VERSION: ", versions[0].Version)
}

var AidokuSource = Source{
	Name:        "Aidoku",
	Identifier:  "com.github.Aidoku",
	Subtitle:    "Free and open source manga reader for iOS and iPadOS",
	Description: "Unofficial Aidoku source for easier updates.",
	Website:     "https://github.com/Aidoku/Aidoku/",
	Apps:        []App{AidokuApp},
}
var booting = false

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	updateApp()
}
