package rules

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/utils"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/secrets"
	"github.com/zricethezav/gitleaks/v8/config"
	"github.com/zricethezav/gitleaks/v8/regexp"
)

func AWS() *config.Rule {
	// define rule
	r := config.Rule{
		RuleID:      "aws-access-token",
		Description: "Identified a pattern that may indicate AWS credentials, risking unauthorized cloud resource access and data breaches on AWS platforms.",
		Regex:       regexp.MustCompile(`\b((?:A3T[A-Z0-9]|AKIA|ASIA|ABIA|ACCA)[A-Z2-7]{16})\b`),
		Entropy:     3,
		Keywords: []string{
			// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-unique-ids
			"A3T",  // todo: might not be a valid AWS token
			"AKIA", // Access key
			"ASIA", // Temporary (AWS STS) access key
			"ABIA", // AWS STS service bearer token
			"ACCA", // Context-specific credential
		},
		Allowlists: []*config.Allowlist{
			{
				Regexes: []*regexp.Regexp{
					regexp.MustCompile(`.+EXAMPLE$`),
				},
			},
		},
	}

	// validate
	tps := utils.GenerateSampleSecrets("AWS", "AKIALALEMEL33243OLIB") // gitleaks:allow
	// current AWS tokens cannot contain [0,1,8,9], so their entropy is slightly lower than expected.
	tps = append(tps, utils.GenerateSampleSecrets("AWS", "AKIA"+secrets.NewSecret("[A-Z2-7]{16}"))...)
	tps = append(tps, utils.GenerateSampleSecrets("AWS", "ASIA"+secrets.NewSecret("[A-Z2-7]{16}"))...)
	tps = append(tps, utils.GenerateSampleSecrets("AWS", "ABIA"+secrets.NewSecret("[A-Z2-7]{16}"))...)
	tps = append(tps, utils.GenerateSampleSecrets("AWS", "ACCA"+secrets.NewSecret("[A-Z2-7]{16}"))...)
	fps := []string{
		`key = AKIAXXXXXXXXXXXXXXXX`,           // Low entropy
		`aws_access_key: AKIAIOSFODNN7EXAMPLE`, // Placeholder
		`msgstr "Näytä asiakirjamallikansio."`, // Lowercase
		`TODAYINASIAASACKOFRICEFELLOVER`,       // wrong length
		`CTTCATAGGGTTCACGCTGTGTAAT-ACG--CCTGAGGC-CACA-AGGGGACTTCAGCAACCGTCGGG-GATTC-ATTGCCA-A--TGGAAGCAATC-TA-TGGGTTA-TCGCGGAGTCCGCAAAGACGGCCAGTATG-AAGCAGATTTCGCAC-CAATGTGACTGCATTTCGTG-ATCGGGGTAAGTA-TC-GCCGATTC-GC--CCGTCCA-AGT-CGAAG-TA--GGCAATATAAAGCTGC-CATTGCCGAAGCTATCTCGCTA-TACTTGAT-AATCGGCGG-TAG-CACAG-GTCGCAGTATCG-AC-T--AGG-CCTCTCAAAAGTT-GGGTCCCGGCCTCTGGGAAAAACACCTCT-A-AGCGTCAATCAGCTCGGTTTCGCATATTA-TGATATCCCCCGTTGACCAATTGA--TAGTACCCGAGCTTACCGTCGG-ATTCTGGAGTCTT-ATGAGGTTACCGACGA-CGCAGTACCATAAGT-GCGCAATTTGACTGTTCCCGTCGAGTAACCA-AGCTTTGCTCA-CCGGGATGCGCGCCGATGTGACCAGGGGGCGCATGTTACATTGAC-A-GCTGGATCATGTTATGAC-GTGGGTC-ATGCTAAAAGCCTAAAGGACGGT-GCATTAGTAT-TACCGGGACCTCATATCAATGCGCTCGCTAGTTCCTCTTCTCTTGATAACGTATATGCGTCAGGCGCCCGTCCGCCTCCAATACGTG-ACAACGTC-AGTACTGAGCCTC--AA-ACATCGTCTTGTTCG-CC-TACAAAGGATCGGTAGAAAACTCAATATTCGGGTATAAGGTCGTAGGAAGTGTGTCGCCCAGGGCCG-CTAGA-AGCGCACACAAGCG-CTCCTGTCAAGGAGTTG-GTGAAAA-ATGAAC--GACT-ATTGCGTCAC--CTACCTCT-AAGTTTTT-GACAATTTCATGGACGAATTGA-AGCGTCCACAAGCATCTGCCGTAGATATGCGGTAGGTTTTTACATATG-TCACTGCAGAGTCACGGACA-CACATCGCTGTCAAAATGCTCGTACCTAGT-GT-TTGCGATCCCCC-GCGGCATTA-TCTTTTGAACCCTCGTCCCTGTGG-CTCTGATGATTGAG-GTCTGTA-TTCCCTCGTTGTGGGGGGATTGGACCTT-TGTATAGGTTCTTTAACCG-ATGGGGGGCCG--ATCGA-A-TA-TGCTCCTGTTTGCCCCGAACCTT-ACCTCGG-TCCAGACA-CTAAGAAAAACCCC-C-ACTGTAAGGTGCTGAGCCTTTGGATAGCC-CGCGAATGAT-CC-TAGTTGACAA-CTGAACGCGCTCGAACA-TGCCC-GCCCTCTGA--CTGCTGTCTG-GCACCTTTAGACACGCGTCGAC-CATATATT-AGCGCTGTCTGTGG-AGGT-TGTGTCTTGTTGCTCA-CT-CATTATCTGT-AACTGGCTCC-CTC-CCAT-TGGCGTCTTTACACCAACCGCTAGGTTACAGTGCA-TCTAGCGCCTATTATCAGGGCGT-TTGCAGCGGCGCGGTGGCTATGT-GTTAGACATATC-CTTACACTGTATGCTAG-AGCAAGCCAC-TCTGAATGGGTTGC-CGATGAATGA-TCTTGATC-GAGCTCGCA-AC---TACATGGAGTCCGAAGTGAACCTACGGATGATCGTATTCCAACACGAGGATC-TATACGTATAGG-A-GGCG-TAATCCACAATTTAGTAACTCTTGACGC---GGATGAAAAT-GTCGTTACACCTTCCAGAGGCTCGG-GTATATATATGACCT--TGTGATTGAGGACGATCTAGAATAA-CT-GT-G-CT-AAAGTACAGTAGTTTCTATGT-GGTAGGTGGAGAATACAGAGTAG-ATGATTC-GTGGGCCACA-C--T-ACTTTCAT-TAGAGCAGAGA-C-GTGAGTGAGTTTTACACTAGCCAGATGGACCG-GTGA-AGTCTAACAGCCACCGCTT-GTGAGGTCGTTTCCCAGTC-ACCCTACTACAGGCAAAAACTCAGTGT-CC-GTGA-GTGCGTTAGTGATATTCCCTAACGGTTAGGTAACT-CATGAATTCA-AT-TAAGCGTGTCC-CGGT-CACGCCCCCATGGGGGCCTTCTTGGGAGG--AGCATCTTAT--AT-GCTCACGTGGTT-GATAGG-A-T-AATACACTTTTAGTCAGTCCATCAATAAC-AAAGGAAC---CAGGTGGTCGCAGATA-TCCCGCTGATATAGCACTGTGTAAACTCAGGTGATA-CTAAGC--GCTCTAAT-ACG-CTTAATGGCAATGCCCAGTTC--ACGACTAGCTTATGAGGCCCAGCTATGGACTGCGGC-GGCATGTCGGC-GATGGTTGCCCTCGCCCTAAATTATGTACGA-T-ACCGCCT-CTTGTTCT-CCGCCCATAGGGT-C--AGCAGGCGATAGACTCCCAGAAATTTCCTCGTCGT-CCGAATAAGACTAACACGACTA-TT-CCTCTAC-GT-G-AA-CTTATCA-CAAATG-GCT-TACC-TAGGTGGTGGCAGATCACTTTCCGGTG-TATTACGAATTGACGCATACCGAC-A-CGC-GCTTGTTGGATAATCGACTCTAACCTCCTCTCTGGCACATGT-GCTGGATTACCTC-TATTTT-TCTCGCTTAG--GGAACG-T-CCTCTGTCGCGTGAG-GTACGTTTCACGGGAG-CGGCTTGTTCATGCCACGTCCATTATCGA-AGTG-C-GTAAGG-A-GAGCCCTA--GACTCTACACGGAAA-TC-AAC-GTAGAAGGCTC-A-CT`,
	}
	return utils.Validate(r, tps, fps)
}
