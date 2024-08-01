// message structured
package structures

type Message struct {
	AgentId       string
	AgentHostname string
	AgentCWD      string
	Commands      []Commands
	File          File
}
