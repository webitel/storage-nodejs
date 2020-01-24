--
-- PostgreSQL database dump
--

-- Dumped from database version 12.0 (Debian 12.0-1.pgdg100+1)
-- Dumped by pg_dump version 12.0 (Debian 12.0-1.pgdg100+1)

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

ALTER TABLE IF EXISTS ONLY storage.file_backend_profiles DROP CONSTRAINT IF EXISTS file_backend_profiles_file_backend_profile_type_id_fk;
DROP INDEX IF EXISTS storage.file_backend_profiles_acl_id_uindex;
ALTER TABLE IF EXISTS ONLY storage.upload_file_jobs DROP CONSTRAINT IF EXISTS upload_file_jobs_pkey;
ALTER TABLE IF EXISTS ONLY storage.session DROP CONSTRAINT IF EXISTS session_pkey;
ALTER TABLE IF EXISTS ONLY storage.schedulers DROP CONSTRAINT IF EXISTS schedulers_pkey;
ALTER TABLE IF EXISTS ONLY storage.remove_file_jobs DROP CONSTRAINT IF EXISTS remove_file_jobs_pkey;
ALTER TABLE IF EXISTS ONLY storage.media_files DROP CONSTRAINT IF EXISTS media_files_pkey;
ALTER TABLE IF EXISTS ONLY storage.jobs DROP CONSTRAINT IF EXISTS jobs_pkey;
ALTER TABLE IF EXISTS ONLY storage.files DROP CONSTRAINT IF EXISTS files_pkey;
ALTER TABLE IF EXISTS ONLY storage.file_backend_profiles DROP CONSTRAINT IF EXISTS file_backend_profiles_pkey;
ALTER TABLE IF EXISTS ONLY storage.file_backend_profiles_acl DROP CONSTRAINT IF EXISTS file_backend_profiles_acl_pk;
ALTER TABLE IF EXISTS ONLY storage.file_backend_profile_type DROP CONSTRAINT IF EXISTS file_backend_profile_type_pkey;
ALTER TABLE IF EXISTS storage.upload_file_jobs ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.schedulers ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.remove_file_jobs ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.media_files ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.files ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.file_backend_profiles_acl ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.file_backend_profiles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS storage.file_backend_profile_type ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS storage.upload_file_jobs_id_seq;
DROP TABLE IF EXISTS storage.upload_file_jobs;
DROP TABLE IF EXISTS storage.session;
DROP SEQUENCE IF EXISTS storage.schedulers_id_seq;
DROP TABLE IF EXISTS storage.schedulers;
DROP SEQUENCE IF EXISTS storage.remove_file_jobs_id_seq;
DROP TABLE IF EXISTS storage.remove_file_jobs;
DROP SEQUENCE IF EXISTS storage.media_files_id_seq;
DROP TABLE IF EXISTS storage.media_files;
DROP TABLE IF EXISTS storage.jobs;
DROP SEQUENCE IF EXISTS storage.files_id_seq;
DROP TABLE IF EXISTS storage.files;
DROP SEQUENCE IF EXISTS storage.file_backend_profiles_id_seq;
DROP SEQUENCE IF EXISTS storage.file_backend_profiles_acl_id_seq;
DROP TABLE IF EXISTS storage.file_backend_profiles_acl;
DROP TABLE IF EXISTS storage.file_backend_profiles;
DROP SEQUENCE IF EXISTS storage.file_backend_profile_type_id_seq;
DROP TABLE IF EXISTS storage.file_backend_profile_type;
DROP SCHEMA IF EXISTS storage;
--
-- Name: storage; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA storage;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: file_backend_profile_type; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.file_backend_profile_type (
    id integer NOT NULL,
    name character varying(50),
    code character varying(10)
);


--
-- Name: file_backend_profile_type_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.file_backend_profile_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: file_backend_profile_type_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.file_backend_profile_type_id_seq OWNED BY storage.file_backend_profile_type.id;


--
-- Name: file_backend_profiles; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.file_backend_profiles (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    expire_day integer DEFAULT 0 NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    disabled boolean,
    max_size_mb integer DEFAULT 0 NOT NULL,
    properties jsonb NOT NULL,
    type_id integer NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint NOT NULL,
    data_size double precision DEFAULT 0 NOT NULL,
    data_count bigint DEFAULT 0 NOT NULL,
    created_by bigint,
    updated_by bigint,
    domain_id bigint,
    description character varying DEFAULT ''::character varying
);


--
-- Name: file_backend_profiles_acl; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.file_backend_profiles_acl (
    id bigint NOT NULL,
    dc bigint NOT NULL,
    grantor bigint NOT NULL,
    subject bigint NOT NULL,
    access smallint DEFAULT 0 NOT NULL,
    object bigint NOT NULL
);


--
-- Name: file_backend_profiles_acl_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.file_backend_profiles_acl_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: file_backend_profiles_acl_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.file_backend_profiles_acl_id_seq OWNED BY storage.file_backend_profiles_acl.id;


--
-- Name: file_backend_profiles_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.file_backend_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: file_backend_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.file_backend_profiles_id_seq OWNED BY storage.file_backend_profiles.id;


--
-- Name: files; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.files (
    id bigint NOT NULL,
    domain character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    size bigint NOT NULL,
    mime_type character varying(20),
    properties jsonb NOT NULL,
    instance character varying(20) NOT NULL,
    uuid character varying(36) NOT NULL,
    profile_id integer,
    created_at bigint,
    removed boolean,
    not_exists boolean
);


--
-- Name: files_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.files_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: files_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.files_id_seq OWNED BY storage.files.id;


--
-- Name: jobs; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.jobs (
    id character varying(26) NOT NULL,
    type character varying(32),
    priority bigint,
    schedule_id bigint,
    schedule_time bigint,
    create_at bigint,
    start_at bigint,
    last_activity_at bigint,
    status character varying(32),
    progress bigint,
    data character varying(1024)
);


--
-- Name: media_files; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.media_files (
    id bigint NOT NULL,
    domain character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    size bigint NOT NULL,
    mime_type character varying(40),
    properties jsonb,
    instance character varying(20),
    created_by text,
    created_at bigint,
    updated_by text,
    updated_at bigint
);


--
-- Name: media_files_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.media_files_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: media_files_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.media_files_id_seq OWNED BY storage.media_files.id;


--
-- Name: remove_file_jobs; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.remove_file_jobs (
    id integer NOT NULL,
    file_id bigint NOT NULL,
    created_at bigint,
    created_by character varying(50)
);


--
-- Name: remove_file_jobs_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.remove_file_jobs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: remove_file_jobs_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.remove_file_jobs_id_seq OWNED BY storage.remove_file_jobs.id;


--
-- Name: schedulers; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.schedulers (
    id bigint NOT NULL,
    cron_expression character varying(50) NOT NULL,
    type character varying(50) NOT NULL,
    name character varying(50) NOT NULL,
    description character varying(500),
    time_zone character varying(50),
    created_at bigint NOT NULL,
    enabled boolean NOT NULL
);


--
-- Name: schedulers_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.schedulers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: schedulers_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.schedulers_id_seq OWNED BY storage.schedulers.id;


--
-- Name: session; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.session (
    key text NOT NULL,
    token character varying(500),
    user_id character varying(26),
    domain character varying(100)
);


--
-- Name: upload_file_jobs; Type: TABLE; Schema: storage; Owner: -
--

CREATE TABLE storage.upload_file_jobs (
    id bigint NOT NULL,
    state integer,
    name character varying(100) NOT NULL,
    uuid character varying(36) NOT NULL,
    mime_type character varying(36),
    size bigint NOT NULL,
    email_msg character varying(500),
    email_sub character varying(150),
    instance character varying(10),
    created_at bigint NOT NULL,
    updated_at bigint,
    attempts integer NOT NULL,
    domain_id bigint NOT NULL
);


--
-- Name: upload_file_jobs_id_seq; Type: SEQUENCE; Schema: storage; Owner: -
--

CREATE SEQUENCE storage.upload_file_jobs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: upload_file_jobs_id_seq; Type: SEQUENCE OWNED BY; Schema: storage; Owner: -
--

ALTER SEQUENCE storage.upload_file_jobs_id_seq OWNED BY storage.upload_file_jobs.id;


--
-- Name: file_backend_profile_type id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profile_type ALTER COLUMN id SET DEFAULT nextval('storage.file_backend_profile_type_id_seq'::regclass);


--
-- Name: file_backend_profiles id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profiles ALTER COLUMN id SET DEFAULT nextval('storage.file_backend_profiles_id_seq'::regclass);


--
-- Name: file_backend_profiles_acl id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profiles_acl ALTER COLUMN id SET DEFAULT nextval('storage.file_backend_profiles_acl_id_seq'::regclass);


--
-- Name: files id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.files ALTER COLUMN id SET DEFAULT nextval('storage.files_id_seq'::regclass);


--
-- Name: media_files id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.media_files ALTER COLUMN id SET DEFAULT nextval('storage.media_files_id_seq'::regclass);


--
-- Name: remove_file_jobs id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.remove_file_jobs ALTER COLUMN id SET DEFAULT nextval('storage.remove_file_jobs_id_seq'::regclass);


--
-- Name: schedulers id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.schedulers ALTER COLUMN id SET DEFAULT nextval('storage.schedulers_id_seq'::regclass);


--
-- Name: upload_file_jobs id; Type: DEFAULT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.upload_file_jobs ALTER COLUMN id SET DEFAULT nextval('storage.upload_file_jobs_id_seq'::regclass);


--
-- Name: file_backend_profile_type file_backend_profile_type_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profile_type
    ADD CONSTRAINT file_backend_profile_type_pkey PRIMARY KEY (id);


--
-- Name: file_backend_profiles_acl file_backend_profiles_acl_pk; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profiles_acl
    ADD CONSTRAINT file_backend_profiles_acl_pk PRIMARY KEY (id);


--
-- Name: file_backend_profiles file_backend_profiles_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profiles
    ADD CONSTRAINT file_backend_profiles_pkey PRIMARY KEY (id);


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.files
    ADD CONSTRAINT files_pkey PRIMARY KEY (id);


--
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);


--
-- Name: media_files media_files_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.media_files
    ADD CONSTRAINT media_files_pkey PRIMARY KEY (id);


--
-- Name: remove_file_jobs remove_file_jobs_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.remove_file_jobs
    ADD CONSTRAINT remove_file_jobs_pkey PRIMARY KEY (id);


--
-- Name: schedulers schedulers_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.schedulers
    ADD CONSTRAINT schedulers_pkey PRIMARY KEY (id);


--
-- Name: session session_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (key);


--
-- Name: upload_file_jobs upload_file_jobs_pkey; Type: CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.upload_file_jobs
    ADD CONSTRAINT upload_file_jobs_pkey PRIMARY KEY (id);


--
-- Name: file_backend_profiles_acl_id_uindex; Type: INDEX; Schema: storage; Owner: -
--

CREATE UNIQUE INDEX file_backend_profiles_acl_id_uindex ON storage.file_backend_profiles_acl USING btree (id);


--
-- Name: file_backend_profiles file_backend_profiles_file_backend_profile_type_id_fk; Type: FK CONSTRAINT; Schema: storage; Owner: -
--

ALTER TABLE ONLY storage.file_backend_profiles
    ADD CONSTRAINT file_backend_profiles_file_backend_profile_type_id_fk FOREIGN KEY (type_id) REFERENCES storage.file_backend_profile_type(id);


--
-- PostgreSQL database dump complete
--

