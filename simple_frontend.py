import streamlit as st
import requests
import time


# Function to send the question to the server
def send_question(question):
    url = "http://localhost:8080/submit_data"
    headers = {"Content-Type": "application/json"}
    data = {"question": question}
    response = requests.post(url, headers=headers, json=data)
    if response.status_code == 200:
        return response.json().get("id")
    else:
        st.error("Failed to send question to the server.")
        return None


# Function to get the answer from the server using long polling
def get_answer(question_id):
    url = f"http://localhost:8080/get_processed_data?id={question_id}"
    while True:
        response = requests.get(url)
        if response.status_code == 200:
            answer = response.json().get("data")
            if answer:
                return answer
        time.sleep(3)


st.title("Chatbot Interface")

question = st.text_input("Ask a question:")
if st.button("Submit"):
    if question:
        with st.spinner("Sending question to the server..."):
            question_id = send_question(question)
        if question_id:
            with st.spinner("Waiting for an answer..."):
                answer = get_answer(question_id)
            st.success("Answer received:")
            st.write(answer)
    else:
        st.warning("Please enter a question.")
