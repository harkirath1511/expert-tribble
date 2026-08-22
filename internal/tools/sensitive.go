package tools

var Sensitive = map[string]string{
	"delete_container":  "permanently removes the container",                                                                                                          
    "kill_container":    "force kills the container process",                                                                                                          
    "stop_container":    "stops a running container",                                                                                                                  
    "delete_image":      "removes a Docker image from disk",                                                                                                           
    "pause_container":   "freezes all processes in the container",                                                                                                     
    "prune_containers":  "removes ALL stopped containers",   
}

func IsSensitive(tool string) (reason string, yes bool){
	reason, yes = Sensitive[tool]
	return
}