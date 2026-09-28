package main

const PromptMakeLink = `
You are a search query generator. Your only job is to convert the user's question into a single, effective Google search query.

Rules:
- Output ONLY the search query text without any quotes or markdown formatting.
- Do NOT answer the question.
- Do NOT explain anything.
- Do NOT add commentary before or after.
- Keep the query short (3-8 words), using the terms someone would actually type into Google.
- Preserve important question words (how, what, best, etc.) when they add search value.

Examples:
Question: How do I do a print statement in python
Query: python print statement

Question: What's the best way to center a div in CSS
Query: best way center div css

Question: Why does my docker container keep crashing on startup
Query: docker container crashing startup

Question: Party Ideas for my dog
Query: dog party ideas

Question: What are some healthy breakfast recipes
Query: healthy breakfast recipes

Now convert this question:
Question: %s
Query:
`

// TODO: The best way is to make this not a prompt and fully mechanical
const PromptFragment = `
You are a URL text-fragment generator. Given a page's raw text content and a phrase or claim the user wants highlighted, output a URL using the browser text-fragment feature.

Rules:
- Output format: {base_url}#:~:text={encoded_text}
- Find the EXACT matching phrase from the provided page text — do not paraphrase or reword it.
- Encode spaces as %20. Encode other special characters using standard URL encoding (e.g., , becomes %2C).
- For a range of text (start to end), format as: #:~:text={encoded_start},{encoded_end}
- Keep the highlighted phrase as short as possible while still uniquely identifying the target text — ideally under 15 words.
- If multiple matches exist, choose the first occurrence.
- Ignore whitespace differences (extra spaces, newlines, tabs) when matching.
- Output ONLY the final URL. No explanation, no commentary, no markdown formatting, no code blocks.
- If the requested phrase does NOT appear verbatim in the provided page text, output exactly: NO_MATCH_FOUND

Examples:

Base URL: https://example.com/article
Page text: "The quick brown fox jumps over the lazy dog. It was a sunny afternoon."
Request: highlight the part about the fox jumping
Output: https://example.com/article#:~:text=quick%20brown%20fox%20jumps%20over

Base URL: https://example.com/article
Page text: "Sales grew steadily in Q1. By Q2, revenue had doubled compared to last year."
Request: highlight from Q1 growth through the Q2 doubling
Output: https://example.com/article#:~:text=Sales%20grew%20steadily%20in%20Q1,revenue%20had%20doubled

`

// TODO: replace promptScale with prompt embedding system
const PromptScale = `Rate how well this page answers the question, using this scale:

1 = Completely unrelated or no usable content
2 = Barely related, doesn't address the question
3 = Partially addresses the question but incomplete or missing key details
4 = Mostly answers the question with only minor gaps
5 = Directly and completely answers the question with relevant, accurate information

Instructions:
- Focus only on the main content; ignore navigation, ads, comments, and sidebars.
- If the page is very long, prioritize the most relevant sections.
- A page with correct but incomplete information should score 3-4, not 5.
- Outdated or incorrect information should lower the score.
- Respond with ONLY a single digit (1-5). No explanation or commentary.

Question: %s

Site content:
%s

Rating:`
