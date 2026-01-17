package main

const PrBodyPrompt = `Generate a plain, concise PR description from these commits:

%s

Guidelines:
- Write naturally, like you're explaining to a teammate
- For simple changes (1-3 small commits), keep it plain - just say what changed and why in 1-2 sentences
- For complex changes (multiple features, refactors), use bullet points to organize
- Match the complexity of the description to the complexity of the change
- Never include meta-commentary like "Here's a PR description" or "This commit does..."
- Never use numbered sections like "1. Brief Summary:" or "2. Key Changes:" - just write the content
- Never use bold headers or excessive formatting - keep it simple
- Avoid filler words: "furthermore", "notably", "arguably", "interestingly", "essentially"
- Avoid buzzwords: "leverage", "streamline", "robust", "innovative", "cutting-edge", "seamless", "empower"
- Avoid vague words: "enhance", "optimize", "transform", "revolutionize", "improve"
- Be specific about what changed, don't use generic praise
- Focus on why the change was made, not how it was implemented

Examples:
Simple change: "Removed unnecessary log statements from admin handlers to reduce log noise."
Complex change: "Added user authentication system:\n- JWT token generation and validation\n- Login/logout endpoints\n- Session management middleware\n- Password hashing with bcrypt"`
