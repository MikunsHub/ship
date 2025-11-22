package main

const PrBodyPrompt = `Analyze these git commits and generate a concise, professional PR description.

Commits:
%s

Please provide:
1. A brief summary (2-3 sentences) of what changed
2. Key changes as bullet points
3. Any notable implementation details

Guidelines:
- Use simple, direct language
- Avoid filler words like "furthermore", "notably", "arguably", "interestingly"
- Don't use buzzwords: "leverage", "streamline", "robust", "innovative", "cutting-edge", "seamless", "empower"
- Skip hedging language: "potentially", "arguably", "could be said that", "in some cases"
- Avoid vague descriptors: "enhance", "optimize", "transform", "revolutionize"
- Be specific about what changed, not generic praise
- Focus on the actual impact and reason for the change
- Write like you're explaining to a teammate, not writing marketing copy

Keep it concise and technical. Focus on the "why" and "what", not the "how".`
