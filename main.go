package main

import (
	"bytes"
		"encoding/json"
			"io"
				"log"
					"net/http"
						"os"

							"github.com/gin-gonic/gin"
								"github.com/joho/godotenv"
								)

								type AIRequest struct {
									Prompt string `json:"prompt"`
									}

									type GeminiRequest struct {
										Contents []struct {
												Parts []struct {
															Text string `json:"text"`
																	} `json:"parts"`
																		} `json:"contents"`
																		}

																		func main() {
																			godotenv.Load()

																				r := gin.Default()

																					r.Use(func(c *gin.Context) {
																							c.Writer.Header().Set("Access-Control-Allow-Origin", "*") 
																									c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
																											c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
																													if c.Request.Method == "OPTIONS" {
																																c.AbortWithStatus(204)
																																			return
																																					}
																																							c.Next()
																																								})

																																									r.POST("/api/ai-assistant", func(c *gin.Context) {
																																											var req AIRequest
																																													if err := c.ShouldBindJSON(&req); err != nil {
																																																c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
																																																			return
																																																					}

																																																							apiKey := os.Getenv("GEMINI_API_KEY")
																																																									if apiKey == "" {
																																																												c.JSON(http.StatusInternalServerError, gin.H{"error": "API key missing"})
																																																															return
																																																																	}

																																																																			systemPrompt := "You are the Coreview AI Assistant. You help users buy building materials and manage their security. Keep answers short and helpful. User asks: " + req.Prompt
																																																																					
																																																																							geminiData := GeminiRequest{
																																																																										Contents: []struct {
																																																																														Parts []struct {
																																																																																			Text string `json:"text"`
																																																																																							} `json:"parts"`
																																																																																										}{
																																																																																														{
																																																																																																			Parts: []struct {
																																																																																																									Text string `json:"text"`
																																																																																																														}{
																																																																																																																				{Text: systemPrompt},
																																																																																																																									},
																																																																																																																													},
																																																																																																																																},
																																																																																																																																		}

																																																																																																																																				jsonData, _ := json.Marshal(geminiData)
																																																																																																																																						url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey
																																																																																																																																								resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
																																																																																																																																										if err != nil {
																																																																																																																																													c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AI"})
																																																																																																																																																return
																																																																																																																																																		}
																																																																																																																																																				defer resp.Body.Close()

																																																																																																																																																						body, _ := io.ReadAll(resp.Body)
																																																																																																																																																								var result map[string]interface{}
																																																																																																																																																										json.Unmarshal(body, &result)

																																																																																																																																																												c.JSON(http.StatusOK, gin.H{
																																																																																																																																																															"reply": result,
																																																																																																																																																																	})
																																																																																																																																																																		})

																																																																																																																																																																			port := os.Getenv("PORT")
																																																																																																																																																																				if port == "" {
																																																																																																																																																																						port = "8080"
																																																																																																																																																																							}
																																																																																																																																																																								r.Run(":" + port)
																																																																																																																																																																								}