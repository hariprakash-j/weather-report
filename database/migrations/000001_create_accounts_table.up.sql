CREATE TABLE scheduler_profile (
	profile_id SERIAL PRIMARY KEY,
	profile_name VARCHAR(20) UNIQUE NOT NULL,
	--- monday
	mon_start TIME NOT NULL,
	mon_stop TIME NOT NULL,
	--- tuesday
	tue_start TIME NOT NULL,
	tue_stop TIME NOT NULL,
	-- wednesday
	wed_start TIME NOT NULL,
	wed_stop TIME NOT NULL,
	-- thursday
	thur_start TIME NOT NULL,
	thur_stop TIME NOT NULL,
	-- friday
	fri_start TIME NOT NULL,
	fri_stop TIME NOT NULL,
	--- saturday
	sat_start TIME NOT NULL,
	sat_stop TIME NOT NULL,
	--- sunday
	sun_start TIME NOT NULL,
	sun_stop TIME NOT NULL
);

CREATE TYPE aws_account_type AS ENUM ('development', 'staging', 'production');

CREATE TABLE aws_account (
	account_id VARCHAR(12) PRIMARY KEY,
	accunnt_type aws_account_type NOT NULL,
	default_profile INT REFERENCES scheduler_profile(profile_id)
);

CREATE TYPE aws_resource_type AS ENUM ('ec2-instance', 'simple-rds-instance');

CREATE TABLE aws_resource (
	resource_aws_id VARCHAR(50) UNIQUE NOT NULL,
	resource_id BIGSERIAL PRIMARY KEY,
	aws_account VARCHAR(12) REFERENCES aws_account(account_id) NOT NULL,
	resource_type aws_resource_type NOT NULL
);
