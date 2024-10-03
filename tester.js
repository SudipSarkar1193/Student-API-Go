const fn = async () => {
    const student = {
        id: Math.floor(Math.random() * 1000), // Random ID between 0 and 999
        fulName: `Student ${Math.floor(Math.random() * 100)}`, // Random name
        email: `student${Math.floor(Math.random() * 100)}@example.com` // Random email
    };
    
    try {
        const response = await fetch("http://127.0.0.1:5000/api/students", {
            method: "POST",
            // headers: {
            //     "Content-Type": "application/json" // Specify that we are sending JSON
            // },
            body: JSON.stringify(student), 
        });

        // Check if the response is OK (status in the range 200-299)
        if (!response.ok) {
            // Try to get the error as JSON, or fallback to plain text
            let errorMessage;
            const contentType = response.headers.get("content-type");
            console.log("------------------------")
            console.log("response.headers => ",response.headers)
            console.log("------------------------")
            if (contentType && contentType.includes("application/json")) {

                const errorData = await response.json();
                console.log("errorData =>",errorData)
                console.log("------------------------")
                errorMessage = errorData.errorMessage || errorData.body ||  "Unknown JSON error";
            } 
            else {

                errorMessage = await response.text(); // Read as plain text
                console.log("response.status =>",response.status)
                console.log("------------------------")
                console.log("errorMessage => ",errorMessage)
                console.log("------------------------")
            }

            console.error(`Error: ${errorMessage} (status: ${response.status})`);
            return;
        }

        // Parse the successful response as JSON
        const data = await response.json();
        
        // Log the parsed data
        console.log("---");
        console.log(data);
        console.log();
       
    } catch (error) {
        // Log any errors that occur during the fetch
        console.error('Error fetching data:', error);
    }
};

fn();
