CREATE TABLE "mc_logs" {
    pkey UUID NOT NULL DEFAULT gen_random_uuid() , 
    CONSTRAINT pkey_mc_logs PRIMARY KEY ( pkey ) ,
    "server" VARCHAR(50) NOT NULL ,
    "server_type" VARCHAR(50) NOT NULL ,
    "timestamp" TIMESTAMP NOT NULL DEFAULT now() ,
    "level" VARCHAR(10) NOT NULL ,
    "source" VARCHAR(50) NOT NULL ,
    "message" TEXT NOT NULL
}
