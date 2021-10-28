# awsclientconfig

This package allows us to create an STS-targeted user for creation of a session on the fly within AWS. It's smart enough
to skip steps when the answers are obvious re: whether we need to do STS assignment or not. This allows us to focus the
"blast-radius" of a user in aws with too many inherent permissions to a set of roles in AWS that are accessible to a 
service user. (least privilege for the task)