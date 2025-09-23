package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func downloadCardImage(cardName string) error {
	baseURL := "https://api.scryfall.com/cards/named"
	params := url.Values{}
	params.Set("format", "image")
	params.Set("version", "normal")
	params.Set("exact", cardName)
	apiURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("card not found: %s", cardName)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	// Determine file extension from Content-Type
	ext := ".jpg"
	ct := resp.Header.Get("Content-Type")
	if ct == "image/png" {
		ext = ".png"
	}
	filename := cardName + ext
	filepath := filepath.Join(ImgDir, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	fmt.Printf("Saved image to %s\n", filepath)
	return nil
}
