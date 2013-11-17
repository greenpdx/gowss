--
-- PostgreSQL database dump
--

-- Dumped from database version 9.1.9
-- Dumped by pg_dump version 9.1.9
-- Started on 2013-09-26 19:30:29 PDT

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- TOC entry 172 (class 3079 OID 11681)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2007 (class 0 OID 0)
-- Dependencies: 172
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- TOC entry 173 (class 3079 OID 16495)
-- Dependencies: 5
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- TOC entry 2008 (class 0 OID 0)
-- Dependencies: 173
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- TOC entry 174 (class 3079 OID 16400)
-- Dependencies: 5
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 2009 (class 0 OID 0)
-- Dependencies: 174
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

--
-- TOC entry 232 (class 1255 OID 16545)
-- Dependencies: 571 5
-- Name: passchg(integer, text, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION passchg(ui integer, op text, ip text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
	rec RECORD;
	dig TEXT;
BEGIN
	SELECT * INTO rec FROM usrpas WHERE "usridx" = ui;
	IF found THEN
		dig = digest(op || rec.salt, 'sha256');
		IF dig = rec.passwd THEN
			dig = digest(ip || rec.salt, 'sha256');
			UPDATE usrpas SET passwd = dig WHERE usridx = ui;
			RETURN TRUE;
		ELSE
			RAISE NOTICE 'BAD password';
		END IF;
	ELSE
		RAISE NOTICE 'User Not Found';
	END IF;
	RETURN FALSE;
END;
$$;


--
-- TOC entry 230 (class 1255 OID 16544)
-- Dependencies: 571 5
-- Name: passchk(integer, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION passchk(ui integer, pas text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$DECLARE
	rec RECORD;
	dig TEXT;
BEGIN
	SELECT * INTO rec FROM usrpas WHERE "usridx" = ui;
	IF found THEN
		dig = digest(pas || rec.salt, 'sha256');
		IF dig = rec.passwd THEN
			RETURN TRUE;
		ELSE
			RAISE NOTICE 'BAD password';
		END IF;
	ELSE
		RAISE NOTICE 'User Not Found';
	END IF;
	RETURN FALSE;
END;$$;


--
-- TOC entry 231 (class 1255 OID 16555)
-- Dependencies: 571 5
-- Name: passchkinfo(text, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION passchkinfo(logi text, pass text) RETURNS record
    LANGUAGE plpgsql
    AS $$BEGIN
END;$$;


--
-- TOC entry 235 (class 1255 OID 16669)
-- Dependencies: 5 571
-- Name: passnew(integer, text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION passnew(ui integer, pass text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$DECLARE
	slt TEXT;
	dig TEXT;
BEGIN
	slt = gen_salt('md5');
	dig = digest(pass || slt, 'sha256');
	INSERT INTO usrpas (usridx,passwd,salt) VALUES (ui,dig,slt);
	RETURN TRUE;
	EXCEPTION
	WHEN unique_violation THEN
	RETURN FALSE;
	WHEN foreign_key_violation THEN
	RETURN FALSE;
END;$$;


--
-- TOC entry 199 (class 1255 OID 16485)
-- Dependencies: 5 571
-- Name: sessins(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION sessins() RETURNS trigger
    LANGUAGE plpgsql
    AS $$DECLARE
	rec RECORD;
BEGIN
IF NEW.svalid IS NULL THEN
	SELECT * INTO rec FROM sess WHERE sesidx = NEW.sesidx AND svalid = '1';
	IF NOT FOUND THEN
		INSERT INTO sess (uid,slast,ip,port,svalid) VALUES (NEW.uid, NEW.start, NEW.ip, NEW.port, TRUE);
	ELSIF OLD.start > now() THEN
		UPDATE sess SET valid = FALSE WHERE sesidx = rec.sesidx;
		INSERT INTO sess (uid,slast,ip,port,svalid) VALUES (NEW.uid, NEW.start, NEW.ip, NEW.port, TRUE);
	ELSE
		UPDATE sess SET slast = now() WHERE sesidx = rec.sesidx;
	END IF;
END IF;
RETURN NEW;
END
$$;


--
-- TOC entry 234 (class 1255 OID 16697)
-- Dependencies: 571 5
-- Name: useradd(character varying, character varying, character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION useradd(plogin character varying, pemail character varying, ppass character varying) RETURNS integer
    LANGUAGE plpgsql
    AS $$DECLARE
	idx INT;
	bol BOOLEAN;
BEGIN
	IF NOT EXISTS (SELECT usridx FROM usrinfo WHERE login = plogin) THEN 
		INSERT INTO usrinfo ("login","email") VALUES (plogin, pemail)
		RETURNING usridx INTO idx;
--		RAISE NOTICE 'IDX %', idx;
		PERFORM passnew(idx,ppass);
		RETURN idx;

	END IF;
	RETURN -1;
	EXCEPTION
	WHEN unique_violation THEN
		RETURN -2;
END;$$;


--
-- TOC entry 233 (class 1255 OID 16671)
-- Dependencies: 571 5
-- Name: usrins(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION usrins() RETURNS trigger
    LANGUAGE plpgsql
    AS $$DECLARE
	rec RECORD;
	pass TEXT;
BEGIN

	SELECT * INTO rec FROM usrinfo WHERE login = NEW.login;
	IF NOT FOUND THEN
		RETURN NULL;
	END IF;

	RAISE NOTICE 'INS UINFO %', NEW;

	RETURN NULL;
--	ELSIF OLD.start > now() THEN
--		UPDATE sess SET valid = FALSE WHERE sesidx = rec.sesidx;
--	ELSE
--		UPDATE sess SET slast = now() WHERE sesidx = rec.sesidx;
--		RETURN NULL;
--	END IF;

END;
$$;


SET default_with_oids = false;

--
-- TOC entry 162 (class 1259 OID 16390)
-- Dependencies: 1972 1973 5
-- Name: gofwconf; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE gofwconf (
    cfgidx integer NOT NULL,
    name character varying(64) NOT NULL,
    rndkey numeric DEFAULT random() NOT NULL,
    uuid uuid DEFAULT uuid_generate_v4() NOT NULL
);


--
-- TOC entry 161 (class 1259 OID 16388)
-- Dependencies: 162 5
-- Name: hzcconf_cfgidx_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE hzcconf_cfgidx_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2010 (class 0 OID 0)
-- Dependencies: 161
-- Name: hzcconf_cfgidx_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE hzcconf_cfgidx_seq OWNED BY gofwconf.cfgidx;


--
-- TOC entry 165 (class 1259 OID 16442)
-- Dependencies: 5
-- Name: nextsession; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE nextsession
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 164 (class 1259 OID 16421)
-- Dependencies: 1975 1976 5
-- Name: qrytst; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE qrytst (
    idx integer NOT NULL,
    dbstring character varying DEFAULT ''::character varying NOT NULL,
    dbint bigint DEFAULT (random() * power((10)::double precision, (15)::double precision)) NOT NULL
);


--
-- TOC entry 163 (class 1259 OID 16419)
-- Dependencies: 5 164
-- Name: qrytst_idx_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE qrytst_idx_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2011 (class 0 OID 0)
-- Dependencies: 163
-- Name: qrytst_idx_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE qrytst_idx_seq OWNED BY qrytst.idx;


--
-- TOC entry 167 (class 1259 OID 16446)
-- Dependencies: 1978 1979 5
-- Name: sess; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE sess (
    sessidx bigint NOT NULL,
    uid integer NOT NULL,
    slast timestamp with time zone,
    ip inet NOT NULL,
    svalid boolean DEFAULT false NOT NULL,
    port integer NOT NULL,
    sstart timestamp with time zone DEFAULT now() NOT NULL,
    intval integer
);


--
-- TOC entry 166 (class 1259 OID 16444)
-- Dependencies: 5 167
-- Name: sess_sesidx_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE sess_sesidx_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2012 (class 0 OID 0)
-- Dependencies: 166
-- Name: sess_sesidx_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE sess_sesidx_seq OWNED BY sess.sessidx;


--
-- TOC entry 169 (class 1259 OID 16546)
-- Dependencies: 1981 1983 5
-- Name: usrinfo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE usrinfo (
    name character varying,
    login character varying NOT NULL,
    email character varying NOT NULL,
    priv smallint[] DEFAULT '{1}'::smallint[] NOT NULL,
    usridx integer NOT NULL,
    pass text DEFAULT '*'::text NOT NULL
);


--
-- TOC entry 171 (class 1259 OID 16572)
-- Dependencies: 5 169
-- Name: usrinfo_usridx_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE usrinfo_usridx_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2013 (class 0 OID 0)
-- Dependencies: 171
-- Name: usrinfo_usridx_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE usrinfo_usridx_seq OWNED BY usrinfo.usridx;


--
-- TOC entry 168 (class 1259 OID 16534)
-- Dependencies: 1980 5
-- Name: usrpas; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE usrpas (
    passwd text NOT NULL,
    salt text DEFAULT gen_salt('md5'::text) NOT NULL,
    usridx integer NOT NULL
);


--
-- TOC entry 170 (class 1259 OID 16561)
-- Dependencies: 5 168
-- Name: usrpas_usridx_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE usrpas_usridx_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2014 (class 0 OID 0)
-- Dependencies: 170
-- Name: usrpas_usridx_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE usrpas_usridx_seq OWNED BY usrpas.usridx;


--
-- TOC entry 1971 (class 2604 OID 16393)
-- Dependencies: 162 161 162
-- Name: cfgidx; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY gofwconf ALTER COLUMN cfgidx SET DEFAULT nextval('hzcconf_cfgidx_seq'::regclass);


--
-- TOC entry 1974 (class 2604 OID 16424)
-- Dependencies: 164 163 164
-- Name: idx; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY qrytst ALTER COLUMN idx SET DEFAULT nextval('qrytst_idx_seq'::regclass);


--
-- TOC entry 1977 (class 2604 OID 16449)
-- Dependencies: 167 166 167
-- Name: sessidx; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY sess ALTER COLUMN sessidx SET DEFAULT nextval('sess_sesidx_seq'::regclass);


--
-- TOC entry 1982 (class 2604 OID 16574)
-- Dependencies: 171 169
-- Name: usridx; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrinfo ALTER COLUMN usridx SET DEFAULT nextval('usrinfo_usridx_seq'::regclass);


--
-- TOC entry 1985 (class 2606 OID 16399)
-- Dependencies: 162 162 2002
-- Name: hzcconf_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY gofwconf
    ADD CONSTRAINT hzcconf_pkey PRIMARY KEY (cfgidx);


--
-- TOC entry 1987 (class 2606 OID 16431)
-- Dependencies: 164 164 2002
-- Name: qrytst_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY qrytst
    ADD CONSTRAINT qrytst_pkey PRIMARY KEY (idx);


--
-- TOC entry 1989 (class 2606 OID 16456)
-- Dependencies: 167 167 2002
-- Name: sess_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY sess
    ADD CONSTRAINT sess_pkey PRIMARY KEY (sessidx);


--
-- TOC entry 1993 (class 2606 OID 16684)
-- Dependencies: 169 169 2002
-- Name: usrinfo_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrinfo
    ADD CONSTRAINT usrinfo_email_key UNIQUE (email);


--
-- TOC entry 1995 (class 2606 OID 16682)
-- Dependencies: 169 169 2002
-- Name: usrinfo_login_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrinfo
    ADD CONSTRAINT usrinfo_login_key UNIQUE (login);


--
-- TOC entry 1997 (class 2606 OID 16623)
-- Dependencies: 169 169 2002
-- Name: usrinfo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrinfo
    ADD CONSTRAINT usrinfo_pkey PRIMARY KEY (usridx);


--
-- TOC entry 1991 (class 2606 OID 16656)
-- Dependencies: 168 168 2002
-- Name: usrpas_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrpas
    ADD CONSTRAINT usrpas_pkey PRIMARY KEY (usridx);


--
-- TOC entry 1999 (class 2620 OID 16486)
-- Dependencies: 199 167 2002
-- Name: sesins; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER sesins BEFORE INSERT ON sess FOR EACH ROW EXECUTE PROCEDURE sessins();


--
-- TOC entry 2000 (class 2620 OID 16672)
-- Dependencies: 233 169 2002
-- Name: usrins; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER usrins BEFORE INSERT ON usrinfo FOR EACH ROW EXECUTE PROCEDURE usrins();

ALTER TABLE usrinfo DISABLE TRIGGER usrins;


--
-- TOC entry 1998 (class 2606 OID 16664)
-- Dependencies: 1996 169 168 2002
-- Name: usrpas_usridx_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY usrpas
    ADD CONSTRAINT usrpas_usridx_fkey FOREIGN KEY (usridx) REFERENCES usrinfo(usridx);


-- Completed on 2013-09-26 19:30:29 PDT

--
-- PostgreSQL database dump complete
--

