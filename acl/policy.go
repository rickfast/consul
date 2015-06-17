package acl

import (
	"fmt"

	"github.com/hashicorp/hcl"
)

const (
	KeyPolicyDeny      = "deny"
	KeyPolicyRead      = "read"
	KeyPolicyWrite     = "write"
	ServicePolicyDeny  = "deny"
	ServicePolicyRead  = "read"
	ServicePolicyWrite = "write"
	EventPolicyRead    = "read"
	EventPolicyFire    = "fire"
	EventPolicyDeny    = "deny"
	ExecPolicyAllow    = "allow"
	ExecPolicyDeny     = "deny"
)

// Policy is used to represent the policy specified by
// an ACL configuration.
type Policy struct {
	ID       string           `hcl:"-"`
	Keys     []*KeyPolicy     `hcl:"key,expand"`
	Services []*ServicePolicy `hcl:"service,expand"`
	Events   []*EventPolicy   `hcl:"event,expand"`
	Exec     []*ExecPolicy    `hcl:"exec,expand"`
}

// KeyPolicy represents a policy for a key
type KeyPolicy struct {
	Prefix string `hcl:",key"`
	Policy string
}

func (k *KeyPolicy) GoString() string {
	return fmt.Sprintf("%#v", *k)
}

// ServicePolicy represents a policy for a service
type ServicePolicy struct {
	Name   string `hcl:",key"`
	Policy string
}

func (k *ServicePolicy) GoString() string {
	return fmt.Sprintf("%#v", *k)
}

// EventPolicy represents a user event policy.
type EventPolicy struct {
	Event  string `hcl:",key"`
	Policy string
}

func (e *EventPolicy) GoString() string {
	return fmt.Sprintf("%#v", *e)
}

// ExecPolicy represents a remote exec policy.
type ExecPolicy struct {
	Command string `hcl:",key"`
	Policy  string
}

func (e *ExecPolicy) GoString() string {
	return fmt.Sprintf("%#v", *e)
}

// Parse is used to parse the specified ACL rules into an
// intermediary set of policies, before being compiled into
// the ACL
func Parse(rules string) (*Policy, error) {
	// Decode the rules
	p := &Policy{}
	if rules == "" {
		// Hot path for empty rules
		return p, nil
	}

	if err := hcl.Decode(p, rules); err != nil {
		return nil, fmt.Errorf("Failed to parse ACL rules: %v", err)
	}

	// Validate the key policy
	for _, kp := range p.Keys {
		switch kp.Policy {
		case KeyPolicyDeny:
		case KeyPolicyRead:
		case KeyPolicyWrite:
		default:
			return nil, fmt.Errorf("Invalid key policy: %#v", kp)
		}
	}

	// Validate the service policy
	for _, sp := range p.Services {
		switch sp.Policy {
		case ServicePolicyDeny:
		case ServicePolicyRead:
		case ServicePolicyWrite:
		default:
			return nil, fmt.Errorf("Invalid service policy: %#v", sp)
		}
	}

	// Validate the user event policies
	for _, ep := range p.Events {
		switch ep.Policy {
		case EventPolicyRead:
		case EventPolicyFire:
		case EventPolicyDeny:
		default:
			return nil, fmt.Errorf("Invalid event policy: %#v", ep)
		}
	}

	// Validate the exec policies
	for _, ep := range p.Exec {
		switch ep.Policy {
		case ExecPolicyAllow:
		case ExecPolicyDeny:
		default:
			return nil, fmt.Errorf("Invalid exec policy: %#v", ep)
		}
	}

	return p, nil
}
