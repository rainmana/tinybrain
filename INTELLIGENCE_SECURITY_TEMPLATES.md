# Intelligence Security Templates for TinyBrain

This document provides pre-configured JSON templates for storing intelligence findings, threat actors, attack campaigns, IOCs, patterns, and correlations in TinyBrain. These templates are designed to support OSINT, HUMINT, SIGINT, and other intelligence gathering activities.

## Intelligence Finding Templates

### OSINT Finding Template
```json
{
  "title": "OSINT Finding: [Target] Social Media Intelligence",
  "content": "Social media analysis reveals [specific findings] about [target]. Key observations include [observations]. Potential attack vectors identified: [vectors]. Confidence level: [high/medium/low].",
  "content_type": "intelligence",
  "category": "intelligence",
  "intelligence_type": "osint",
  "classification": "unclassified",
  "threat_level": "medium",
  "geographic_scope": "national",
  "attribution": "unknown",
  "ioc_type": "domain",
  "ioc_value": "example.com",
  "mitre_tactic": "TA0043",
  "mitre_technique": "T1591",
  "mitre_procedure": "T1591.001",
  "kill_chain_phase": "reconnaissance",
  "risk_score": 6.5,
  "impact_score": 7.0,
  "likelihood_score": 6.0,
  "confidence": 0.8,
  "priority": 7,
  "tags": ["osint", "social-media", "reconnaissance", "target-analysis"],
  "source": "Social Media Platforms",
  "metadata": {
    "platforms_analyzed": ["Twitter", "LinkedIn", "Facebook"],
    "time_range": "2024-01-01 to 2024-01-31",
    "analysis_tools": ["Maltego", "theHarvester"],
    "key_personnel_identified": ["John Doe", "Jane Smith"],
    "technical_details": "Found references to internal systems and technologies"
  }
}
```

### HUMINT Finding Template
```json
{
  "title": "HUMINT Finding: [Source] Intelligence Report",
  "content": "Human intelligence source reports [specific information] regarding [target/activity]. Source reliability: [high/medium/low]. Information corroborated by [sources]. Actionable intelligence: [actions].",
  "content_type": "intelligence",
  "category": "intelligence",
  "intelligence_type": "humint",
  "classification": "confidential",
  "threat_level": "high",
  "geographic_scope": "regional",
  "attribution": "Source Alpha",
  "ioc_type": "email",
  "ioc_value": "suspicious@example.com",
  "mitre_tactic": "TA0001",
  "mitre_technique": "T1566",
  "mitre_procedure": "T1566.001",
  "kill_chain_phase": "delivery",
  "risk_score": 8.5,
  "impact_score": 9.0,
  "likelihood_score": 8.0,
  "confidence": 0.9,
  "priority": 9,
  "tags": ["humint", "source-intelligence", "threat-actor", "campaign"],
  "source": "Human Source",
  "metadata": {
    "source_reliability": "high",
    "information_currency": "recent",
    "corroboration_sources": ["OSINT", "SIGINT"],
    "actionable_items": ["immediate_threat", "planned_attack"],
    "sensitivity_level": "high"
  }
}
```

### SIGINT Finding Template
```json
{
  "title": "SIGINT Finding: [Target] Communications Intelligence",
  "content": "Signals intelligence analysis reveals [specific findings] from [target communications]. Technical details: [details]. Threat indicators: [indicators]. Recommended actions: [actions].",
  "content_type": "intelligence",
  "category": "intelligence",
  "intelligence_type": "sigint",
  "classification": "secret",
  "threat_level": "critical",
  "geographic_scope": "international",
  "attribution": "APT Group",
  "ioc_type": "ip",
  "ioc_value": "192.168.1.100",
  "mitre_tactic": "TA0011",
  "mitre_technique": "T1071",
  "mitre_procedure": "T1071.001",
  "kill_chain_phase": "c2",
  "risk_score": 9.5,
  "impact_score": 9.5,
  "likelihood_score": 9.0,
  "confidence": 0.95,
  "priority": 10,
  "tags": ["sigint", "communications", "apt", "c2", "infrastructure"],
  "source": "Signals Intelligence",
  "metadata": {
    "communication_types": ["email", "chat", "voip"],
    "encryption_status": "encrypted",
    "protocols_identified": ["HTTPS", "SSH", "RDP"],
    "geographic_origins": ["Russia", "China"],
    "technical_indicators": ["custom_protocols", "steganography"]
  }
}
```

## Threat Actor Templates

### APT Group Template
```json
{
  "name": "APT[Number]: [Group Name]",
  "aliases": ["Alias1", "Alias2", "Alias3"],
  "description": "Advanced Persistent Threat group known for [capabilities and activities]. Active since [timeframe]. Primary targets: [targets]. Attribution: [attribution].",
  "motivation": "Financial gain / Espionage / Disruption",
  "capabilities": ["Spear phishing", "Zero-day exploits", "Custom malware", "Lateral movement"],
  "targets": ["Government", "Financial", "Healthcare", "Critical Infrastructure"],
  "tools": ["Custom malware", "Commercial tools", "Open source tools"],
  "techniques": ["T1566.001", "T1055", "T1071.001", "T1021"],
  "attribution": "High confidence attribution to [country/organization]",
  "confidence": 0.9,
  "threat_level": "critical",
  "geographic_scope": "international",
  "metadata": {
    "first_observed": "2020-01-01",
    "last_observed": "2024-01-31",
    "campaigns": ["Campaign1", "Campaign2"],
    "iocs": ["hash1", "domain1", "ip1"],
    "tactics_used": ["TA0001", "TA0002", "TA0003"],
    "victim_countries": ["US", "UK", "DE"],
    "sophistication_level": "high"
  }
}
```

### Cybercriminal Group Template
```json
{
  "name": "[Group Name] Cybercriminal Organization",
  "aliases": ["Alias1", "Alias2"],
  "description": "Cybercriminal group specializing in [crime types]. Known for [modus operandi]. Estimated size: [size]. Primary motivation: [motivation].",
  "motivation": "Financial gain",
  "capabilities": ["Ransomware", "Banking trojans", "Carding", "Money laundering"],
  "targets": ["Small businesses", "Individuals", "Financial institutions"],
  "tools": ["Ransomware-as-a-Service", "Cryptocurrency", "Dark web marketplaces"],
  "techniques": ["T1486", "T1055", "T1071.001"],
  "attribution": "Suspected to be based in [country]",
  "confidence": 0.7,
  "threat_level": "high",
  "geographic_scope": "international",
  "metadata": {
    "estimated_revenue": "$1M+",
    "preferred_payment": "Bitcoin",
    "ransom_amounts": ["$10K", "$50K", "$100K"],
    "victim_countries": ["US", "UK", "CA"],
    "arrests": ["Name1", "Name2"],
    "law_enforcement_actions": ["Operation1", "Operation2"]
  }
}
```

## Attack Campaign Templates

### APT Campaign Template
```json
{
  "name": "Operation [Name]",
  "description": "Long-term APT campaign targeting [sectors/regions]. Duration: [start] to [end]. Attribution: [attribution]. Impact: [impact].",
  "threat_actors": ["APT1", "APT2"],
  "targets": ["Government", "Critical Infrastructure", "Financial"],
  "techniques": ["T1566.001", "T1055", "T1071.001", "T1021", "T1486"],
  "tools": ["Custom malware", "Living off the land", "Commercial tools"],
  "iocs": ["hash1", "domain1", "ip1", "email1"],
  "start_date": "2023-01-01T00:00:00Z",
  "end_date": "2023-12-31T23:59:59Z",
  "status": "completed",
  "threat_level": "critical",
  "geographic_scope": "international",
  "confidence": 0.95,
  "metadata": {
    "victim_countries": ["US", "UK", "DE", "FR"],
    "sectors_targeted": ["Government", "Energy", "Financial"],
    "data_exfiltrated": ["PII", "Intellectual Property", "Credentials"],
    "estimated_damage": "$10M+",
    "law_enforcement_response": ["Arrests", "Sanctions", "Indictments"],
    "attribution_confidence": "high"
  }
}
```

### Ransomware Campaign Template
```json
{
  "name": "[Ransomware Name] Campaign [Year]",
  "description": "Ransomware campaign using [ransomware family]. Targets: [targets]. Distribution: [method]. Impact: [impact].",
  "threat_actors": ["Ransomware Group"],
  "targets": ["Healthcare", "Education", "Small Business"],
  "techniques": ["T1486", "T1055", "T1071.001"],
  "tools": ["[Ransomware Name]", "RDP", "Email"],
  "iocs": ["hash1", "domain1", "ip1"],
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": null,
  "status": "active",
  "threat_level": "high",
  "geographic_scope": "international",
  "confidence": 0.8,
  "metadata": {
    "ransomware_family": "[Name]",
    "ransom_amounts": ["$5K", "$25K", "$100K"],
    "payment_methods": ["Bitcoin", "Monero"],
    "victim_countries": ["US", "UK", "CA", "AU"],
    "sectors_affected": ["Healthcare", "Education", "Government"],
    "total_victims": 1000,
    "total_ransom_paid": "$5M+"
  }
}
```

## IOC Templates

### Domain IOC Template
```json
{
  "type": "domain",
  "value": "malicious.example.com",
  "description": "Malicious domain used for [purpose]. First seen: [date]. Last seen: [date]. Associated with: [campaign/actor].",
  "threat_level": "high",
  "confidence": 0.9,
  "first_seen": "2024-01-01T00:00:00Z",
  "last_seen": "2024-01-31T23:59:59Z",
  "source": "Threat Intelligence Feed",
  "attribution": "APT Group",
  "campaigns": ["Campaign1", "Campaign2"],
  "techniques": ["T1071.001", "T1566.001"],
  "tags": ["malware", "c2", "phishing"],
  "metadata": {
    "dns_records": ["A", "AAAA", "MX"],
    "ip_addresses": ["1.2.3.4", "5.6.7.8"],
    "subdomains": ["sub1", "sub2"],
    "ssl_certificate": "self-signed",
    "whois_data": "privacy_protected",
    "reputation": "malicious"
  }
}
```

### IP Address IOC Template
```json
{
  "type": "ip",
  "value": "192.168.1.100",
  "description": "Malicious IP address used for [purpose]. Geolocation: [country]. ISP: [provider]. Associated with: [campaign/actor].",
  "threat_level": "medium",
  "confidence": 0.8,
  "first_seen": "2024-01-01T00:00:00Z",
  "last_seen": "2024-01-31T23:59:59Z",
  "source": "Network Monitoring",
  "attribution": "Unknown",
  "campaigns": ["Campaign1"],
  "techniques": ["T1071.001"],
  "tags": ["c2", "scanning", "brute_force"],
  "metadata": {
    "geolocation": {
      "country": "Russia",
      "region": "Moscow",
      "city": "Moscow",
      "latitude": 55.7558,
      "longitude": 37.6176
    },
    "isp": "Example ISP",
    "asn": "AS12345",
    "ports_open": [22, 80, 443, 3389],
    "services": ["SSH", "HTTP", "HTTPS", "RDP"],
    "reputation": "suspicious"
  }
}
```

### File Hash IOC Template
```json
{
  "type": "hash",
  "value": "a1b2c3d4e5f6789012345678901234567890abcd",
  "description": "Malicious file hash. File type: [type]. Size: [size]. Associated with: [campaign/actor].",
  "threat_level": "critical",
  "confidence": 0.95,
  "first_seen": "2024-01-01T00:00:00Z",
  "last_seen": "2024-01-31T23:59:59Z",
  "source": "Malware Analysis",
  "attribution": "APT Group",
  "campaigns": ["Campaign1"],
  "techniques": ["T1055", "T1486"],
  "tags": ["malware", "trojan", "backdoor"],
  "metadata": {
    "file_type": "PE32",
    "file_size": 1024000,
    "hash_algorithm": "SHA256",
    "malware_family": "Trojan.Example",
    "detection_names": ["Malware/Example", "Trojan.Generic"],
    "sandbox_analysis": "available",
    "yara_rules": ["rule1", "rule2"],
    "behavior_indicators": ["network_communication", "file_creation", "registry_modification"]
  }
}
```

## Pattern Templates

### Attack Pattern Template
```json
{
  "name": "Spear Phishing with Malicious Attachment",
  "description": "Attack pattern involving targeted phishing emails with malicious attachments. Common in [sectors]. Success rate: [rate].",
  "pattern_type": "attack",
  "category": "phishing",
  "severity": "high",
  "confidence": 0.9,
  "frequency": 50,
  "examples": [
    "Email with malicious Word document",
    "PDF with embedded JavaScript",
    "Excel file with macros"
  ],
  "mitigations": [
    "Email filtering and scanning",
    "User awareness training",
    "Application whitelisting",
    "Sandboxing of attachments"
  ],
  "detections": [
    "Email security gateways",
    "Endpoint detection and response",
    "User reporting",
    "Behavioral analysis"
  ],
  "tags": ["phishing", "malware", "social-engineering"],
  "metadata": {
    "mitre_techniques": ["T1566.001", "T1566.002"],
    "target_sectors": ["Government", "Financial", "Healthcare"],
    "common_attachments": [".doc", ".pdf", ".xls"],
    "success_factors": ["urgency", "authority", "familiarity"],
    "time_patterns": ["business_hours", "end_of_week"]
  }
}
```

### Behavioral Pattern Template
```json
{
  "name": "Lateral Movement via RDP",
  "description": "Behavioral pattern of attackers using RDP for lateral movement within compromised networks. Common after initial compromise.",
  "pattern_type": "behavioral",
  "category": "lateral_movement",
  "severity": "medium",
  "confidence": 0.8,
  "frequency": 30,
  "examples": [
    "RDP connections between workstations",
    "Unusual login times",
    "Multiple failed login attempts",
    "Privilege escalation via RDP"
  ],
  "mitigations": [
    "Network segmentation",
    "RDP restrictions",
    "Multi-factor authentication",
    "Privileged access management"
  ],
  "detections": [
    "Network monitoring",
    "Authentication logs",
    "Behavioral analytics",
    "Privilege monitoring"
  ],
  "tags": ["lateral_movement", "rdp", "privilege_escalation"],
  "metadata": {
    "mitre_techniques": ["T1021.001"],
    "common_tools": ["mstsc", "rdesktop", "freerdp"],
    "network_indicators": ["port_3389", "rdp_traffic"],
    "log_indicators": ["event_id_4624", "event_id_4625"],
    "time_patterns": ["after_hours", "weekends"]
  }
}
```

## Correlation Templates

### Temporal Correlation Template
```json
{
  "source_finding_id": "finding_1",
  "target_finding_id": "finding_2",
  "correlation_type": "temporal",
  "strength": 0.9,
  "confidence": 0.8,
  "evidence": "Both findings occurred within 24 hours and share similar IOCs",
  "description": "Temporal correlation between phishing campaign and malware deployment",
  "weight": 1.0,
  "direction": "unidirectional",
  "metadata": {
    "time_difference": "2 hours",
    "shared_iocs": ["domain1", "ip1"],
    "attack_sequence": ["phishing", "malware_deployment"],
    "confidence_factors": ["timing", "ioc_overlap", "attack_chain"]
  }
}
```

### Logical Correlation Template
```json
{
  "source_finding_id": "finding_1",
  "target_finding_id": "finding_2",
  "correlation_type": "logical",
  "strength": 0.95,
  "confidence": 0.9,
  "evidence": "Both findings are part of the same attack campaign based on TTPs and attribution",
  "description": "Logical correlation between reconnaissance and exploitation phases",
  "weight": 1.5,
  "direction": "bidirectional",
  "metadata": {
    "shared_ttps": ["T1591", "T1566.001"],
    "attribution_match": "APT Group",
    "campaign_indicators": ["similar_tools", "target_overlap"],
    "confidence_factors": ["ttp_similarity", "attribution", "timing"]
  }
}
```

## MITRE ATT&CK Integration Templates

### Tactic Mapping Template
```json
{
  "tactic_id": "TA0001",
  "tactic_name": "Initial Access",
  "description": "The adversary is trying to get into your network.",
  "techniques": [
    {
      "technique_id": "T1566",
      "technique_name": "Phishing",
      "description": "Adversaries may send phishing messages to gain access to victim systems.",
      "subtechniques": [
        {
          "subtechnique_id": "T1566.001",
          "subtechnique_name": "Spearphishing Attachment",
          "description": "Adversaries may send spearphishing emails with a malicious attachment."
        }
      ]
    }
  ],
  "mitigations": [
    "M1056: Pre-boot Authentication",
    "M1026: Privileged Account Management"
  ],
  "detections": [
    "DS0017: Command",
    "DS0022: File"
  ]
}
```

### Technique Mapping Template
```json
{
  "technique_id": "T1566.001",
  "technique_name": "Spearphishing Attachment",
  "tactic": "TA0001",
  "description": "Adversaries may send spearphishing emails with a malicious attachment.",
  "platforms": ["Office 365", "G Suite", "Outlook"],
  "permissions_required": ["User"],
  "data_sources": ["Email", "File"],
  "mitigations": [
    "M1056: Pre-boot Authentication",
    "M1026: Privileged Account Management"
  ],
  "detections": [
    "DS0017: Command",
    "DS0022: File"
  ],
  "examples": [
    "Malicious Word document with macros",
    "PDF with embedded JavaScript",
    "Excel file with malicious formulas"
  ]
}
```

## Usage Guidelines

### Template Selection
1. **Intelligence Type**: Choose the appropriate intelligence type (OSINT, HUMINT, SIGINT, etc.)
2. **Classification Level**: Set appropriate classification based on sensitivity
3. **Threat Level**: Assess and set threat level based on impact and likelihood
4. **Geographic Scope**: Define the geographic scope of the intelligence

### Metadata Best Practices
1. **Source Attribution**: Always include source information and reliability
2. **Temporal Data**: Include first seen, last seen, and analysis timeframes
3. **Technical Details**: Include relevant technical indicators and analysis
4. **Confidence Scoring**: Use consistent confidence scoring methodology

### Correlation Guidelines
1. **Evidence-Based**: Base correlations on concrete evidence
2. **Confidence Levels**: Set appropriate confidence levels for correlations
3. **Direction**: Specify whether correlations are unidirectional or bidirectional
4. **Weight**: Use weights to indicate correlation importance

### MITRE ATT&CK Integration
1. **Tactic Mapping**: Map findings to appropriate MITRE tactics
2. **Technique Identification**: Identify specific techniques and sub-techniques
3. **Procedure Details**: Include procedure-level details when available
4. **Kill Chain**: Map to appropriate kill chain phases

These templates provide a comprehensive foundation for intelligence gathering and analysis within TinyBrain, enabling structured storage and analysis of security intelligence data.
