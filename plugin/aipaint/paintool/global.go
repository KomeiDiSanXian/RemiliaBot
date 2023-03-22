package paintool

var (
	paintURL = "http://ipv4.rinkore.com:7866"
	txt2img  = "/sdapi/v1/txt2img"
)

// samplers
const (
	samplerEulerA      = "Euler a"
	samplerEuler       = "Euler"
	samplerLMS         = "LMS"
	samplerHeun        = "Heun"
	samplerDPM2        = "DPM2"
	samplerDPM2A       = "DPM2 a"
	samplerDPMPP2SA    = "DPM++ 2S a"
	samplerDPMPP2M     = "DPM++ 2M"
	samplerDPMPPSDE    = "DPM++ SDE"
	samplerDPMFast     = "DPM fast"
	samplerDPMAdaptive = "DPM adaptive"
	samplerLMSK        = "LMS karras"
	samplerDPM2K       = "DPM2 karras"
	samplerDPM2AK      = "DPM2 a karras"
	samplerDPMPP2SK    = "DPM++ 2S a karras"
	samplerDPMPP2MK    = "DPM++ 2M Karras"
	samplerDPMPPSDEk   = "DPM++ SDE karras"
	samplerDDIM        = "DDIM"
	samplerPLMS        = "PLMS"
)

// negative prompts
const (
	DefaultNegtivePrompt = "lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, lowres, text, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, nsfw, lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry, missing arms, long neck, humpbacked"
)
