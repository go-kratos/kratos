package chinese

var (
	hk2s = `{
		"name": "Traditional Chinese (Hong Kong standard) to Simplified Chinese",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "TSPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "HKVariantsRevPhrases.txt"
				}, {
					"type": "txt",
					"file": "HKVariantsRev.txt"
				}]
			}
		}, {
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TSPhrases.txt"
				}, {
					"type": "txt",
					"file": "TSCharacters.txt"
				}]
			}
		}]
	}
	`
	s2hk = `{
		"name": "Simplified Chinese to Traditional Chinese (Hong Kong standard)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "STPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "STPhrases.txt"
				}, {
					"type": "txt",
					"file": "STCharacters.txt"
				}]
			}
		}, {
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "HKVariantsPhrases.txt"
				}, {
					"type": "txt",
					"file": "HKVariants.txt"
				}]
			}
		}]
	}
	`
	s2t = `{
		"name": "Simplified Chinese to Traditional Chinese",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "STPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "STPhrases.txt"
				}, {
					"type": "txt",
					"file": "STCharacters.txt"
				}]
			}
		}]
	}
	`
	s2tw = `{
		"name": "Simplified Chinese to Traditional Chinese (Taiwan standard)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "STPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "STPhrases.txt"
				}, {
					"type": "txt",
					"file": "STCharacters.txt"
				}]
			}
		}, {
			"dict": {
				"type": "txt",
				"file": "TWVariants.txt"
			}
		}]
	}
	`
	s2twp = `{
		"name": "Simplified Chinese to Traditional Chinese (Taiwan standard, with phrases)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "STPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "STPhrases.txt"
				}, {
					"type": "txt",
					"file": "STCharacters.txt"
				}]
			}
		}, {
			"dict": {
				"type": "txt",
				"file": "TWPhrases.txt"
			}
		}, {
			"dict": {
				"type": "txt",
				"file": "TWVariants.txt"
			}
		}]
	}
	`
	t2hk = `{
		"name": "Traditional Chinese to Traditional Chinese (Hong Kong standard)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "HKVariants.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "txt",
				"file": "HKVariants.txt"
			}
		}]
	}
	`
	t2s = `{
		"name": "Traditional Chinese to Simplified Chinese",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "TSPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TSPhrases.txt"
				}, {
					"type": "txt",
					"file": "TSCharacters.txt"
				}]
			}
		}]
	}
	`
	t2tw = `{
		"name": "Traditional Chinese to Traditional Chinese (Taiwan standard)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "TWVariants.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "txt",
				"file": "TWVariants.txt"
			}
		}]
	}
	`
	tw2s = `{
		"name": "Traditional Chinese (Taiwan standard) to Simplified Chinese",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "TSPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TWVariantsRevPhrases.txt"
				}, {
					"type": "txt",
					"file": "TWVariantsRev.txt"
				}]
			}
		}, {
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TSPhrases.txt"
				}, {
					"type": "txt",
					"file": "TSCharacters.txt"
				}]
			}
		}]
	}
	`
	tw2sp = `{
		"name": "Traditional Chinese (Taiwan standard) to Simplified Chinese (with phrases)",
		"segmentation": {
			"type": "mmseg",
			"dict": {
				"type": "txt",
				"file": "TSPhrases.txt"
			}
		},
		"conversion_chain": [{
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TWVariantsRevPhrases.txt"
				}, {
					"type": "txt",
					"file": "TWVariantsRev.txt"
				}]
			}
		}, {
			"dict": {
				"type": "txt",
				"file": "TWPhrasesRev.txt"
			}
		}, {
			"dict": {
				"type": "group",
				"dicts": [{
					"type": "txt",
					"file": "TSPhrases.txt"
				}, {
					"type": "txt",
					"file": "TSCharacters.txt"
				}]
			}
		}]
	}
	`
)
