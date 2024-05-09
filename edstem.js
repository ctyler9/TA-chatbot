var URL = "http://localhost:8080"

// Function to make a GET request to the server
async function fetchData(id) {
    const getUrl = `${URL}/get_processed_data?id=${id}`;
    const resp = await fetch(getUrl);
    return await resp.json();
}

// Long polling function
async function longPolling(id) {
    while (true) {
        try {
            // Make a GET request to fetch data
            const data = await fetchData(id);
			//
            // Process the data here, for example, log it
            console.log('Received data:', data);

            // Check if there's any new data
            if (!lastData || JSON.stringify(data) !== JSON.stringify(lastData)) {
                // New data detected
                lastData = data;
                return data;
            }

            // Wait for a certain period before making the next request
            await new Promise(resolve => setTimeout(resolve, 5000)); // Wait for 5 seconds
        } catch (error) {
            console.error('Error:', error);
        }
    }
}

edBot.postComment(async (comment, actions) => {
    if (comment.plaintext.includes('LLMBOT')) {
        var submissionQuestion = comment.plaintext;
        try {
            // POST JSON to an external API and receive JSON response
            const resp = await fetch(`${URL}/submit_data`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    question: submissionQuestion,
                }),
            });
            const body = await resp.json();

            // Assuming the JSON response contains an ID
            const lookupID = body.id;

            // Start long polling for data
            const submissionAnswer = await longPolling(lookupID);

            // Once new data is received, continue with actions
            actions.comment(submissionAnswer.data);
            
            // Assuming the GET response contains the processed data
            console.log('Processed data:', submissionAnswer.data);


        } catch (error) {
            console.error('Error:', error);
            actions.comment("Error connecting to server/timeout. Please try again")
        }
    }
});
