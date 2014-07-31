for v in `seq 1900 2000`;do ./fetch -url "http://jiraserver.com/" -user jira-user -password p -issue DCE-$v; done

