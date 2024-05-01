BEGIN;

CREATE TYPE cat_race AS ENUM (
    'Persian',
	'Maine Coon',
	'Siamese',
	'Ragdoll',
	'Bengal',
	'Sphynx',
	'British Shorthair',
	'Abyssinian',
	'Scottish Fold',
	'Birman'
    );


COMMIT;