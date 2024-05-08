# TA Chatbot
The goal for this project is to assist students for any class with LLMs + class material

## Design 
As LLMs have an increasingly large context size, the need for RAG is probably diminishing. With this in mind, the idea for this project is to have a simple key, value map (provided as a json file) where the server will fetch the relevant parts from and send to the LLM to process. 

The process works as follows: you tell students to put in their request a specific keyword for the question they are asking about. For example: "syllabus", "HW1Q1", "Project". The server will find this specific keyword and fetch the relevant values set in the json. With only the relevant documents now in hand, it sends a direct request to the LLM to service the students question. No messy/annoying vector embeddings/extra complexity. 

### Example Scenario

1. **Student Asks Question**: 
   - A student asks a question related to their homework assignment. The question must include a specific keyword, like "HW3Q1."

2. **Keyword Parsing**: 
   - The system parses the student's question to identify the keyword, such as "HW3Q1."

3. **Document Lookup (JSON)**: 
   - The system looks up relevant documents associated with the keyword "HW3Q1" in a JSON file containing mappings of keywords to specific documents.

4. **Build Prompt**: 
   - The system constructs a custom prompt using the prompt path, the student's question, and the relevant documents retrieved from the lookup step.

5. **Prompt Sent to Language Model (LLM)**: 
   - The constructed prompt is then sent to a Language Model (LLM) to generate an answer to the student's question based on the provided context.

This process ensures that the student's question is appropriately contextualized, allowing for more accurate and relevant responses from the Language Model.

## Installation
Since bindings for Golang use sub-modules, need to clone repo as mentoined in instructions, build, and reference it: https://github.com/go-skynet/go-llama.cpp/tree/master

## To run
go build .\
./main (or go run .) -d {docs_path} -m {model_path} -p {prompt_path}\






