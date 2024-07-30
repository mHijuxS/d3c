# Steps for the agent
1. Starts execution
2. Creates unique ID
3. Presents itself to the server
    - Sends id
    - Must have a transmitted packet structured
    - Received and sent message must have a standard structure
4. Wait an x amount of time
5. Connects to the server again to verify if there are any commands to be executed locally
6. Executes commands locally
7. Send response to server

# Steps for the server
1. Listens to any inbound connection
2. Receives the ID from the agent
3. Checks if ID is already registered
4. Check if there are any response from previous commands
5. Sends message to agent
6. Receives response from agent

# Create Server CLI
1. select <AgentID>
2. show agents
3. Type commands
    - Stored on the message to be sent to the selected agent