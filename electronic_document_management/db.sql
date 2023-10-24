CREATE TABLE "User" (
	"id" serial NOT NULL,
	"email" VARCHAR(255) NOT NULL UNIQUE,
	"full_name" VARCHAR(255) NOT NULL,
	"phone" VARCHAR(255),
	"image" VARCHAR(255),
	"role" VARCHAR(255) NOT NULL,
	CONSTRAINT "User_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Role" (
	"name" VARCHAR(255) NOT NULL,
	CONSTRAINT "Role_pk" PRIMARY KEY ("name")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Request" (
	"id" serial NOT NULL,
	"title" VARCHAR(255) NOT NULL,
	"client" integer NOT NULL,
	"worker" integer NOT NULL,
	"client_object" integer NOT NULL,
	"equipment" integer NOT NULL,
	"contract" integer,
	"description" TEXT,
	"priority" VARCHAR(255) NOT NULL,
	"start_date" VARCHAR(255) NOT NULL,
	"end_date" VARCHAR(255),
	"files" VARCHAR(255),
	"status" VARCHAR(255) NOT NULL,
	CONSTRAINT "Request_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Client" (
	"id" serial NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"inn"	VARCHAR(255),
	"kpp" VARCHAR(255),
	"ogrn" VARCHAR(255),
	"owner" VARCHAR(255),
	"phone" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"address" VARCHAR(255) NOT NULL,
	"create_date" integer NOT NULL,
	"status" BOOLEAN NOT NULL DEFAULT 'true',
	CONSTRAINT "Client_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Objects" (
	"id" serial NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"address" VARCHAR(255) NOT NULL,
	"work_schedule" VARCHAR(255),
	CONSTRAINT "Objects_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Client_objects" (
	"id" serial NOT NULL,
	"client" integer NOT NULL,
	"object" integer NOT NULL,
	CONSTRAINT "Client_objects_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Equipment" (
	"id" serial NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"type" VARCHAR(255) NOT NULL,
	"manufacturer" VARCHAR(255),
	"model" VARCHAR(255) NOT NULL,
	"unique_number" VARCHAR(255),
	"contract" VARCHAR(255),
	"create_date" VARCHAR(255) NOT NULL,
	CONSTRAINT "Equipment_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Client_equipment" (
	"id" serial NOT NULL,
	"equipment" integer NOT NULL,
	"client" integer NOT NULL,
	"object" integer NOT NULL,
	"contract" integer NOT NULL,
	CONSTRAINT "Client_equipment_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Contract" (
	"id" serial NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"client" integer NOT NULL,
	"start_date" VARCHAR(255) NOT NULL,
	"end_date" VARCHAR(255) NOT NULL,
	"amount" integer NOT NULL,
	"file" VARCHAR(255) NOT NULL,
	"status" BOOLEAN NOT NULL DEFAULT 'true',
	CONSTRAINT "Contract_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Request_priority" (
	"name" VARCHAR(255) NOT NULL,
	CONSTRAINT "Request_priority_pk" PRIMARY KEY ("name")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Request_status" (
	"name" VARCHAR(255) NOT NULL,
	"color" VARCHAR(255) NOT NULL,
	CONSTRAINT "Request_status_pk" PRIMARY KEY ("name")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Services" (
	"id" serial NOT NULL,
	"equipment" integer NOT NULL,
	"type" VARCHAR(255) NOT NULL,
	"cost" integer NOT NULL,
	CONSTRAINT "Services_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Services_type" (
	"name" VARCHAR(255) NOT NULL,
	CONSTRAINT "Services_type_pk" PRIMARY KEY ("name")
) WITH (
  OIDS=FALSE
);



ALTER TABLE "User" ADD CONSTRAINT "User_fk0" FOREIGN KEY ("role") REFERENCES "Role"("name");


ALTER TABLE "Request" ADD CONSTRAINT "Request_fk0" FOREIGN KEY ("client") REFERENCES "Client"("id");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk1" FOREIGN KEY ("worker") REFERENCES "User"("id");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk2" FOREIGN KEY ("client_object") REFERENCES "Client_objects"("id");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk3" FOREIGN KEY ("equipment") REFERENCES "Equipment"("id");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk4" FOREIGN KEY ("contract") REFERENCES "Contract"("id");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk5" FOREIGN KEY ("priority") REFERENCES "Request_priority"("name");
ALTER TABLE "Request" ADD CONSTRAINT "Request_fk6" FOREIGN KEY ("status") REFERENCES "Request_status"("name");



ALTER TABLE "Client_objects" ADD CONSTRAINT "Client_objects_fk0" FOREIGN KEY ("client") REFERENCES "Client"("id");
ALTER TABLE "Client_objects" ADD CONSTRAINT "Client_objects_fk1" FOREIGN KEY ("object") REFERENCES "Objects"("id");


ALTER TABLE "Client_equipment" ADD CONSTRAINT "Client_equipment_fk0" FOREIGN KEY ("equipment") REFERENCES "Equipment"("id");
ALTER TABLE "Client_equipment" ADD CONSTRAINT "Client_equipment_fk1" FOREIGN KEY ("client") REFERENCES "Client"("id");
ALTER TABLE "Client_equipment" ADD CONSTRAINT "Client_equipment_fk2" FOREIGN KEY ("object") REFERENCES "Objects"("id");
ALTER TABLE "Client_equipment" ADD CONSTRAINT "Client_equipment_fk3" FOREIGN KEY ("contract") REFERENCES "Contract"("id");

ALTER TABLE "Contract" ADD CONSTRAINT "Contract_fk0" FOREIGN KEY ("client") REFERENCES "Client"("id");



ALTER TABLE "Services" ADD CONSTRAINT "Services_fk0" FOREIGN KEY ("equipment") REFERENCES "Equipment"("id");
ALTER TABLE "Services" ADD CONSTRAINT "Services_fk1" FOREIGN KEY ("type") REFERENCES "Services_type"("name");
