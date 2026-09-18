# Liberate - An LLM as a Librarian Rather Than a Teacher

A librarian does not know every single detail, but can guide you in
the right direction so you can learn from the books around you.

Liberate's goal is to "Liberate" people by providing sources of information
rather than hallucinating information and asserting its correctness.
**With Liberate you are the arbiter of the information you receive.**

A locally hosted model with a search engine. Creates the query for
a search using the questions you provide and validates the relevance
of the information to your question. **Bringing control back to the
user since software should always be made for users.**

Why not just use a search engine? Simply put, you don't know
what you don't know, and using the LLM's context we can simplify the
mechanical effort of searching for information to explicitly find
things relevant to your query.

# Roadmap

**In no particular order**

- Reduce time to access links to be comparable to just using a search engine
- [URI fragmentation](https://developer.mozilla.org/en-US/docs/Web/URI/Reference/Fragment) for improved searching, i.e.,
  finding the exact location of relevant information using highlighting.
- Frontend customization and depth
- Flexible LLM options for a wider range of hardware options
- History and recall of queries and responses
- Docker containerization to simplify deployment to one command
- Concurrent LLM querying

# Installation

Note: Liberate currently requires `ollama serve` and an instance of `searxng` running on separate ports,
said available ports can be modified for the `.env`

Liberate currently requires `qwen2.5:3b` to be downloaded and available in your ollama instance. `ollama pull qwen2.5:3b`
or `ollama run qwen2.5:3b` to both test and pull

searxng should have at least this in `settings.yml` since Liberate
interfaces with searxng with json only

```
search:
  formats:
    - json
```

1. `git clone https://github.com/David-Bosnic/Liberate`
2. `cp example.env .env` then modify the ports at which you have ollama and searxng running. Included is
   the port at which the frontend is rendered; change if required.
3. `go build` and `./liberate`
4. Should be running on `http://localhost:8081/` unless .env `FRONTEND_PORT=8081` is modified
