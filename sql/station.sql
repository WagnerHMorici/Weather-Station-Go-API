CREATE TABLE public.estacao (
	id serial4 NOT NULL,
	estacao varchar(100) NULL,
	cidade varchar(100) NULL,
	coordenadas varchar(50) NULL,
	inicio_de_operacao date NULL,
	fim_de_operacao date NULL,
	em_uso bool NULL,
	CONSTRAINT estacao_pkey PRIMARY KEY (id)
);