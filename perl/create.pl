use strict;
use JSON;
use Akamai::Edgegrid;
use LWP::ConsoleLogger::Easy qw( debug_ua );
use Data::Dumper;
use Getopt::Long;

my $section = "default";
my $edgerc = "$ENV{HOME}/.edgerc";
my $debug = 0;
my $version = "v0.1";

GetOptions (
	'config:s' => \$edgerc,
	'section:s' => \$section,
	'property:s' => \$section,
	'debug!' => \$debug,
	'version' => sub {
		print "$0 $version\n";
		exit 0;
	},
) or usage();

my $configname = $ARGV[0];
usage() if not $configname;

my $agent = new Akamai::Edgegrid(
	config_file => "$ENV{HOME}/.edgerc",
	section => $section,
);
my $baseurl = "https://" . $agent->{host};

debug_ua($agent) if $debug;

# Get groups
print "Getting groups\n";
my $req = HTTP::Request->new(GET => $baseurl . '/papi/v1/groups');
$req->content_type('application/json');

my $resp = $agent->request($req);

my $groups_json = decode_json($resp->content);
my %groups;
foreach my $entry (@{$groups_json->{groups}{items}}) {
	$groups{$entry->{groupId}} = $entry->{groupName};
}
die "Couldn't get groups" if not %groups;

# Get property
print "Getting property\n";
my %json;
$json{propertyName} = $configname;
my $req = HTTP::Request->new(POST => $baseurl . '/papi/v1/search/find-by-value');
$req->content_type('application/json');
$req->content(encode_json(\%json));

my $resp = $agent->request($req);

my $content = decode_json($resp->content);
die "Couldn't find property" if not $content;

# Find latest version
my $latestVersion = 0;
foreach my $version (@{$content->{versions}{items}}) {
	$latestVersion = $version->{propertyVersion} if $latestVersion < $version->{propertyVersion};
}

my $groupId = $content->{versions}{items}[0]->{groupId};
my $contractId = $content->{versions}{items}[0]->{contractId};
my $propertyId = $content->{versions}{items}[0]->{propertyId};
my $propertyName = $content->{versions}{items}[0]->{propertyName};

# get a version
print "Getting latest version\n";
my $req = HTTP::Request->new(GET => $baseurl . "/papi/v1/properties/$propertyId/versions/$latestVersion?contractId=$contractId&groupId=$groupId");
$req->content_type('application/json');

my $resp = $agent->request($req);
my $content = decode_json($resp->content);
die "Couldn't find property version" if not $content;

my $productId = $content->{versions}{items}[0]->{productId};
my $ruleFormat = $content->{versions}{items}[0]->{ruleFormat};

# get a rule tree
print "Getting rules\n";
my $req = HTTP::Request->new(GET => $baseurl . "/papi/v1/properties/$propertyId/versions/$latestVersion/rules?contractId=$contractId&groupId=$groupId");
$req->content_type('application/json');

my $resp = $agent->request($req);
my $content = decode_json($resp->content);
die "Couldn't find rules" if not $content;
my $rules = $resp->content;

my $isSecure = $content->{rules}{options}{is_secure} == 1?"true":"false";
my $cpcode;
# search for cpcode in default rule
foreach my $behaviour (@{$content->{rules}{behaviors}}) {
	if ($behaviour->{name} eq "cpCode") {
		$cpcode = $behaviour->{options}{value}{id};
		last;
	}
}

# Get CPCode Name
print "Getting cpcode\n";
my $req = HTTP::Request->new(GET => $baseurl . "/papi/v1/cpcodes/$cpcode?contractId=$contractId&groupId=$groupId");
$req->content_type('application/json');

my $resp = $agent->request($req);
my $content = decode_json($resp->content);
die "Couldn't find cpcode" if not $content;
my $cpcodeName = $content->{cpcodes}{items}[0]->{cpcodeName};

# Get hostnames
print "Getting hostnames\n";
my $req = HTTP::Request->new(GET => $baseurl . "/papi/v1/properties/$propertyId/versions/$latestVersion/hostnames?contractId=$contractId&groupId=$groupId");
$req->content_type('application/json');

my $resp = $agent->request($req);
my $content = decode_json($resp->content);
die "Couldn't find hostnames" if not $content;
my %hostnames;
my %edgehostnames;
foreach my $hostname (@{$content->{hostnames}{items}}) {
	$hostnames{$hostname->{cnameFrom}} = $hostname->{cnameTo};
	$edgehostnames{$hostname->{cnameTo}} = 1;
}


print "Writing configs\n";
open(my $fh, '>', 'rules.json') or die "Can't open rules.json for writing: $!";
print $fh $rules;
close $fh;

open(my $fh, '>', 'property.tf') or die "Can't open property.tf for writing: $!";

#header
print $fh "provider \"akamai\" {\n";
print $fh "  edgerc = \"~/.edgerc\"\n";
print $fh "  papi_section = \"papi\"\n";
print $fh "}\n";
print $fh "\n";

#group
print $fh "data \"akamai_group\" \"group\" {\n";
print $fh "  name = \"" . $groups{$groupId} . "\"\n";
print $fh "}\n";
print $fh "\n";

#contract
print $fh "data \"akamai_contract\" \"contract\" {\n";
print $fh "  group = data.akamai_group.group.name\n";
print $fh "}\n";
print $fh "\n";

#rules
print $fh "data \"template_file\" \"property_json\" {\n";
print $fh "   template = file(\"\${path.module}/rules.json\")\n";
print $fh "}\n";
print $fh "\n";

#cpcode
print $fh "resource \"akamai_cp_code\" \"default-cp-code\" {\n";
print $fh "    product  = \"$productId\"\n";
print $fh "    contract = data.akamai_contract.contract.id\n";
print $fh "    group = data.akamai_group.group.id\n";
print $fh "    name = \"$cpcodeName\"\n";
print $fh "}\n";
print $fh "\n";

#hostnames
foreach my $edgehostname (keys %edgehostnames) {
	my $name = $edgehostname;
	$name =~ s/\./-/g;
	print $fh "resource \"akamai_edge_hostname\" \"$name\" {\n";
	print $fh "    product  = \"$productId\"\n";
	print $fh "    contract = data.akamai_contract.contract.id\n";
	print $fh "    group = data.akamai_group.group.id\n";
	print $fh "    edge_hostname = \"$edgehostname\"\n";
	print $fh "}\n";
	print $fh "\n";
}

my $resourceName = $propertyName;
$resourceName =~ s/\./-/g;

print $fh "resource \"akamai_property\" \"$resourceName\" {\n";
print $fh "  name        = \"$propertyName\"\n";
print $fh "  cp_code     = akamai_cp_code.default-cp-code.id\n";
print $fh "  contact     = [\"\"]\n";
print $fh "  contract = data.akamai_contract.contract.id\n";
print $fh "  group = data.akamai_group.group.id\n";
print $fh "  product     = \"$productId\"\n";
print $fh "  rule_format = \"$ruleFormat\"\n";
print $fh "\n";
print $fh " hostnames    = {\n";

foreach my $hostname (keys %hostnames) {
	my $edgehostname = $hostnames{$hostname};
	$edgehostname =~ s/\./-/g;
	print $fh "       \"$hostname\" = \"akamai_edge_hostname.$edgehostname\",\n";
}
print $fh "  }\n";
print $fh "\n";
print $fh "  rules = data.template_file.property_json.rendered\n";
print $fh "  is_secure = $isSecure\n";
print $fh "}\n";

close $fh;
print "done\n";


sub usage {
	print "Usage: $0 <options> <property>\n";
	print "\n";
	print "Options:\n";
	print "\t--config <config>\tLocation of .edgerc (default: ~/.edgerc)\n";
	print "\t--section <section>\tSection to use within .edgerc\n";
	print "\t--debug\t\tturn on debugging\n";
	print "\t--version\t\tShow version number\n";
	print "Copyright (C) Akamai Technologies, Inc\n";
	print "\n";
	exit(0);
}

