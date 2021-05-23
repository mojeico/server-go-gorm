--
-- PostgreSQL database dump
--

-- Dumped from database version 13.2
-- Dumped by pg_dump version 13.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: test; Type: DATABASE; Schema: -; Owner: postgres
--

-- CREATE DATABASE test WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_US.UTF-8';


ALTER DATABASE test OWNER TO postgres;

\connect test

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: test; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA IF NOT EXISTS test;


ALTER SCHEMA test OWNER TO postgres;

--
-- Name: user_id; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id
    OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: people; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.people
(
    id         bigint NOT NULL,
    first_name text,
    last_name  text,
    age        integer,
    email      text
);


ALTER TABLE test.people
    OWNER TO postgres;

--
-- Name: people_id_seq; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.people_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test.people_id_seq
    OWNER TO postgres;

--
-- Name: people_id_seq; Type: SEQUENCE OWNED BY; Schema: test; Owner: postgres
--

ALTER SEQUENCE test.people_id_seq OWNED BY test.people.id;


--
-- Name: videos; Type: TABLE; Schema: test; Owner: postgres
--

CREATE TABLE test.videos
(
    id          bigint NOT NULL,
    title       text,
    description text,
    url         text,
    person_id   bigint,
    created_ad  timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE test.videos
    OWNER TO postgres;

--
-- Name: videos_id_seq; Type: SEQUENCE; Schema: test; Owner: postgres
--

CREATE SEQUENCE test.videos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test.videos_id_seq
    OWNER TO postgres;

--
-- Name: videos_id_seq; Type: SEQUENCE OWNED BY; Schema: test; Owner: postgres
--

ALTER SEQUENCE test.videos_id_seq OWNED BY test.videos.id;


--
-- Name: people id; Type: DEFAULT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.people
    ALTER COLUMN id SET DEFAULT nextval('test.people_id_seq'::regclass);


--
-- Name: videos id; Type: DEFAULT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.videos
    ALTER COLUMN id SET DEFAULT nextval('test.videos_id_seq'::regclass);


--
-- Data for Name: people; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.people (id, first_name, last_name, age, email) FROM stdin;
1	gleb	mojeico	23	g.mojeico@gmail.com
2	gleb	mojeico	23	g.mojeico@gmail.com
\.


--
-- Data for Name: videos; Type: TABLE DATA; Schema: test; Owner: postgres
--

COPY test.videos (id, title, description, url, person_id, created_ad, updated_at) FROM stdin;
1	video1	video 123	123 url	1	2021-05-21 11:53:20.494934+03	2021-05-21 11:53:20.504907+03
2	video1	video 123	123 url	2	2021-05-21 11:56:09.255831+03	2021-05-21 11:56:09.257118+03
\.


--
-- Name: user_id; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id', 28, true);


--
-- Name: people_id_seq; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.people_id_seq', 2, true);


--
-- Name: videos_id_seq; Type: SEQUENCE SET; Schema: test; Owner: postgres
--

SELECT pg_catalog.setval('test.videos_id_seq', 2, true);


--
-- Name: people people_pkey; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.people
    ADD CONSTRAINT people_pkey PRIMARY KEY (id);


--
-- Name: videos videos_pkey; Type: CONSTRAINT; Schema: test; Owner: postgres
--

ALTER TABLE ONLY test.videos
    ADD CONSTRAINT videos_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

