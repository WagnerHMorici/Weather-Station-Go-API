CREATE TABLE public.registros_estacoes (
	id serial4 NOT NULL,
	temperatura numeric(5, 2) NULL,
	umidade numeric(5, 2) NULL,
	datahora timestamp NULL,
	estacao_fk int4 NULL,
	CONSTRAINT registros_estacoes_pkey PRIMARY KEY (id),
	CONSTRAINT registros_estacoes_estacao_fk_fkey FOREIGN KEY (estacao_fk) REFERENCES public.estacao(id)
);