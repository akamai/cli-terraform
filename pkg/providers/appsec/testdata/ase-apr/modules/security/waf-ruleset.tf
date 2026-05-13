resource "akamai_appsec_waf_ruleset" "default_policy" {
  config_id          = local.config_id
  security_policy_id = akamai_appsec_security_policy.default_policy.security_policy_id
  rules = [
    {
      rule_id     = "950002"
      rule_action = "alert"
    },
    {
      rule_id     = "950006"
      rule_action = "alert"
    },
    {
      rule_id     = "950007"
      rule_action = "alert"
    },
    {
      rule_id     = "950011"
      rule_action = "alert"
    },
    {
      rule_id     = "950118"
      rule_action = "alert"
    },
    {
      rule_id     = "950203"
      rule_action = "alert"
    },
    {
      rule_id     = "950216"
      rule_action = "alert"
    },
    {
      rule_id     = "950220"
      rule_action = "alert"
    },
    {
      rule_id     = "950901"
      rule_action = "alert"
    },
    {
      rule_id     = "951910"
      rule_action = "alert"
    },
    {
      rule_id     = "958003"
      rule_action = "alert"
    },
    {
      rule_id     = "958008"
      rule_action = "alert"
    },
    {
      rule_id     = "958023"
      rule_action = "alert"
    },
    {
      rule_id     = "958034"
      rule_action = "alert"
    },
    {
      rule_id     = "958051"
      rule_action = "alert"
    },
    {
      rule_id     = "958052"
      rule_action = "alert"
    },
    {
      rule_id     = "959070"
      rule_action = "alert"
    },
    {
      rule_id     = "959073"
      rule_action = "alert"
    },
    {
      rule_id     = "959976"
      rule_action = "alert"
    },
    {
      rule_id     = "959977"
      rule_action = "alert"
    },
    {
      rule_id     = "961011"
      rule_action = "alert"
    },
    {
      rule_id     = "961012"
      rule_action = "alert"
    },
    {
      rule_id     = "961034"
      rule_action = "alert"
    },
    {
      rule_id     = "961904"
      rule_action = "alert"
    },
    {
      rule_id     = "961912"
      rule_action = "alert"
    },
    {
      rule_id     = "968230"
      rule_action = "alert"
    },
    {
      rule_id     = "969151"
      rule_action = "alert"
    },
    {
      rule_id     = "970003"
      rule_action = "alert"
    },
    {
      rule_id     = "970004"
      rule_action = "alert"
    },
    {
      rule_id     = "970009"
      rule_action = "alert"
    },
    {
      rule_id     = "970011"
      rule_action = "alert"
    },
    {
      rule_id     = "970013"
      rule_action = "alert"
    },
    {
      rule_id     = "970014"
      rule_action = "alert"
    },
    {
      rule_id     = "970015"
      rule_action = "alert"
    },
    {
      rule_id     = "970118"
      rule_action = "alert"
    },
    {
      rule_id     = "970901"
      rule_action = "alert"
    },
    {
      rule_id     = "970902"
      rule_action = "alert"
    },
    {
      rule_id     = "970903"
      rule_action = "alert"
    },
    {
      rule_id     = "970904"
      rule_action = "alert"
    },
    {
      rule_id     = "973305"
      rule_action = "alert"
    },
    {
      rule_id     = "973307"
      rule_action = "alert"
    },
    {
      rule_id     = "973311"
      rule_action = "alert"
    },
    {
      rule_id     = "973312"
      rule_action = "alert"
    },
    {
      rule_id     = "973335"
      rule_action = "alert"
    },
    {
      rule_id     = "981240"
      rule_action = "alert"
    },
    {
      rule_id     = "981242"
      rule_action = "alert"
    },
    {
      rule_id     = "981243"
      rule_action = "alert"
    },
    {
      rule_id     = "981244"
      rule_action = "alert"
    },
    {
      rule_id     = "981247"
      rule_action = "alert"
    },
    {
      rule_id     = "981248"
      rule_action = "alert"
    },
    {
      rule_id     = "981251"
      rule_action = "alert"
    },
    {
      rule_id     = "981252"
      rule_action = "alert"
    },
    {
      rule_id     = "981253"
      rule_action = "alert"
    },
    {
      rule_id     = "981254"
      rule_action = "alert"
    },
    {
      rule_id     = "981255"
      rule_action = "alert"
    },
    {
      rule_id     = "981256"
      rule_action = "alert"
    },
    {
      rule_id     = "981260"
      rule_action = "alert"
    },
    {
      rule_id     = "981270"
      rule_action = "alert"
    },
    {
      rule_id     = "981276"
      rule_action = "alert"
    },
    {
      rule_id     = "981300"
      rule_action = "alert"
    },
    {
      rule_id     = "981320"
      rule_action = "alert"
    },
    {
      rule_id     = "999002"
      rule_action = "alert"
    },
    {
      rule_id     = "999901"
      rule_action = "alert"
    },
    {
      rule_id     = "999902"
      rule_action = "alert"
    },
    {
      rule_id     = "3000000"
      rule_action = "alert"
    },
    {
      rule_id     = "3000004"
      rule_action = "alert"
    },
    {
      rule_id     = "3000005"
      rule_action = "alert"
    },
    {
      rule_id     = "3000006"
      rule_action = "alert"
    },
    {
      rule_id     = "3000007"
      rule_action = "alert"
    },
    {
      rule_id     = "3000012"
      rule_action = "alert"
    },
    {
      rule_id     = "3000013"
      rule_action = "alert"
    },
    {
      rule_id     = "3000014"
      rule_action = "alert"
    },
    {
      rule_id     = "3000015"
      rule_action = "alert"
    },
    {
      rule_id     = "3000017"
      rule_action = "alert"
    },
    {
      rule_id     = "3000022"
      rule_action = "alert"
    },
    {
      rule_id     = "3000023"
      rule_action = "alert"
    },
    {
      rule_id     = "3000025"
      rule_action = "alert"
    },
    {
      rule_id     = "3000029"
      rule_action = "alert"
    },
    {
      rule_id     = "3000033"
      rule_action = "alert"
    },
    {
      rule_id     = "3000034"
      rule_action = "alert"
    },
    {
      rule_id     = "3000037"
      rule_action = "alert"
    },
    {
      rule_id     = "3000038"
      rule_action = "alert"
    },
    {
      rule_id     = "3000039"
      rule_action = "alert"
    },
    {
      rule_id     = "3000041"
      rule_action = "alert"
    },
    {
      rule_id     = "3000056"
      rule_action = "alert"
    },
    {
      rule_id     = "3000057"
      rule_action = "alert"
    },
    {
      rule_id     = "3000058"
      rule_action = "alert"
    },
    {
      rule_id     = "3000061"
      rule_action = "alert"
    },
    {
      rule_id     = "3000065"
      rule_action = "alert"
    },
    {
      rule_id     = "3000072"
      rule_action = "alert"
    },
    {
      rule_id     = "3000080"
      rule_action = "alert"
    },
    {
      rule_id     = "3000081"
      rule_action = "alert"
    },
    {
      rule_id     = "3000100"
      rule_action = "alert"
    },
    {
      rule_id     = "3000101"
      rule_action = "alert"
    },
    {
      rule_id     = "3000102"
      rule_action = "alert"
    },
    {
      rule_id     = "3000108"
      rule_action = "alert"
    },
    {
      rule_id     = "3000109"
      rule_action = "alert"
    },
    {
      rule_id     = "3000110"
      rule_action = "alert"
    },
    {
      rule_id     = "3000111"
      rule_action = "alert"
    },
    {
      rule_id     = "3000112"
      rule_action = "alert"
    },
    {
      rule_id     = "3000113"
      rule_action = "alert"
    },
    {
      rule_id     = "3000114"
      rule_action = "alert"
    },
    {
      rule_id     = "3000115"
      rule_action = "alert"
    },
    {
      rule_id     = "3000116"
      rule_action = "alert"
    },
    {
      rule_id     = "3000117"
      rule_action = "alert"
    },
    {
      rule_id     = "3000118"
      rule_action = "alert"
    },
    {
      rule_id     = "3000119"
      rule_action = "alert"
    },
    {
      rule_id     = "3000120"
      rule_action = "alert"
    },
    {
      rule_id     = "3000121"
      rule_action = "alert"
    },
    {
      rule_id     = "3000122"
      rule_action = "alert"
    },
    {
      rule_id     = "3000123"
      rule_action = "alert"
    },
    {
      rule_id     = "3000124"
      rule_action = "alert"
    },
    {
      rule_id     = "3000125"
      rule_action = "alert"
    },
    {
      rule_id     = "3000126"
      rule_action = "alert"
    },
    {
      rule_id     = "3000127"
      rule_action = "alert"
    },
    {
      rule_id     = "3000128"
      rule_action = "alert"
    },
    {
      rule_id     = "3000129"
      rule_action = "alert"
    },
    {
      rule_id     = "3000130"
      rule_action = "alert"
    },
    {
      rule_id     = "3000131"
      rule_action = "alert"
    },
    {
      rule_id     = "3000132"
      rule_action = "alert"
    },
    {
      rule_id     = "3000133"
      rule_action = "alert"
    },
    {
      rule_id     = "3000134"
      rule_action = "alert"
    },
    {
      rule_id     = "3000135"
      rule_action = "alert"
    },
    {
      rule_id     = "3000136"
      rule_action = "alert"
    },
    {
      rule_id     = "3000137"
      rule_action = "alert"
    },
    {
      rule_id     = "3000138"
      rule_action = "alert"
    },
    {
      rule_id     = "3000139"
      rule_action = "alert"
    },
    {
      rule_id     = "3000140"
      rule_action = "alert"
    },
    {
      rule_id     = "3000141"
      rule_action = "alert"
    },
    {
      rule_id     = "3000142"
      rule_action = "alert"
    },
    {
      rule_id     = "3000143"
      rule_action = "alert"
    },
    {
      rule_id     = "3000144"
      rule_action = "alert"
    },
    {
      rule_id     = "3000145"
      rule_action = "alert"
    },
    {
      rule_id     = "3000146"
      rule_action = "alert"
    },
    {
      rule_id     = "3000147"
      rule_action = "alert"
    },
    {
      rule_id     = "3000148"
      rule_action = "alert"
    },
    {
      rule_id     = "3000149"
      rule_action = "alert"
    },
    {
      rule_id     = "3000150"
      rule_action = "alert"
    },
    {
      rule_id     = "3000151"
      rule_action = "alert"
    },
    {
      rule_id     = "3000152"
      rule_action = "alert"
    },
    {
      rule_id     = "3000153"
      rule_action = "alert"
    },
    {
      rule_id     = "3000154"
      rule_action = "alert"
    },
    {
      rule_id     = "3000155"
      rule_action = "alert"
    },
    {
      rule_id     = "3000156"
      rule_action = "alert"
    },
    {
      rule_id     = "3000157"
      rule_action = "alert"
    },
    {
      rule_id     = "3000160"
      rule_action = "alert"
    },
    {
      rule_id     = "3000162"
      rule_action = "alert"
    },
    {
      rule_id     = "3000164"
      rule_action = "alert"
    },
    {
      rule_id     = "3000166"
      rule_action = "alert"
    },
    {
      rule_id     = "3000167"
      rule_action = "alert"
    },
    {
      rule_id     = "3000168"
      rule_action = "alert"
    },
    {
      rule_id     = "3000171"
      rule_action = "alert"
    },
    {
      rule_id     = "3000173"
      rule_action = "alert"
    },
    {
      rule_id     = "3000174"
      rule_action = "alert"
    },
    {
      rule_id     = "3000175"
      rule_action = "alert"
    },
    {
      rule_id     = "3000176"
      rule_action = "alert"
    },
    {
      rule_id     = "3000177"
      rule_action = "alert"
    },
    {
      rule_id     = "3000179"
      rule_action = "alert"
    },
    {
      rule_id     = "3000180"
      rule_action = "alert"
    },
    {
      rule_id     = "3000183"
      rule_action = "alert"
    },
    {
      rule_id     = "3000184"
      rule_action = "alert"
    },
    {
      rule_id     = "3000185"
      rule_action = "alert"
    },
  ]
  attack_groups = [
    {
      attack_group        = "POLICY"
      attack_group_action = "deny"
    },
    {
      attack_group        = "WAT"
      attack_group_action = "deny"
    },
    {
      attack_group        = "PROTOCOL"
      attack_group_action = "deny"
    },
    {
      attack_group        = "SQL"
      attack_group_action = "deny"
    },
    {
      attack_group        = "XSS"
      attack_group_action = "deny"
    },
    {
      attack_group        = "CMD"
      attack_group_action = "deny"
    },
    {
      attack_group        = "LFI"
      attack_group_action = "deny"
    },
    {
      attack_group        = "RFI"
      attack_group_action = "deny"
    },
    {
      attack_group        = "PLATFORM"
      attack_group_action = "deny"
    },
  ]
}

