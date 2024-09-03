package asciiArt

func BannerFile(option string) string {
	switch option {
	case "standard":
		return "banners/standard.txt"
	case "shadow":
		return "banners/shadow.txt"
	case "thinkertoy":
		return "banners/thinkertoy.txt"
	default:
		return "invalid bannerfile name"
	}
}
