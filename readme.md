# TA Chatbot
The goal for this project is to assist students for any class with LLMs + class material

## Design 
As LLMs have an increasingly large context size, the need for RAG is probably diminishing. With this in mind, the idea for this project is to have a simple key, value map (provided as a json file) where the server will fetch the relevant parts from and send to the LLM to process. 

The process works as follows: you tell students to put in their request a specific keyword for the question they are asking about. For example: "syllabus", "HW1Q1", "Project". The server will find this specific keyword and fetch the relevant values set in the json. With only the relevant documents now in hand, it sends a direct request to the LLM to service the students question. No messy/annoying vector embeddings/extra complexity. 

## Installation
Since bindings for Golang use sub-modules, need to clone repo as mentoined in instructions, build, and reference it: https://github.com/go-skynet/go-llama.cpp/tree/master




