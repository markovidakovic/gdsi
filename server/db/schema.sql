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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: gender; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.gender AS ENUM (
    'male',
    'female'
);


--
-- Name: handedness; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.handedness AS ENUM (
    'right',
    'left',
    'ambidextrous'
);


--
-- Name: role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.role AS ENUM (
    'developer',
    'admin',
    'user'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.account (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(250) NOT NULL,
    email character varying(250) NOT NULL,
    dob date NOT NULL,
    gender public.gender NOT NULL,
    phone_number character varying(250) NOT NULL,
    password text NOT NULL,
    role public.role DEFAULT 'user'::public.role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: court; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.court (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(250) NOT NULL,
    creator_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: league; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.league (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(250) NOT NULL,
    description character varying(500),
    season_id uuid NOT NULL,
    creator_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: match; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.match (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    court_id uuid NOT NULL,
    scheduled_at timestamp with time zone NOT NULL,
    player_one_id uuid NOT NULL,
    player_two_id uuid NOT NULL,
    winner_id uuid,
    score text,
    season_id uuid NOT NULL,
    league_id uuid NOT NULL,
    creator_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: player; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.player (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    height numeric,
    weight numeric,
    handedness public.handedness,
    racket character varying(250),
    matches_expected integer DEFAULT 0 NOT NULL,
    matches_played integer DEFAULT 0 NOT NULL,
    matches_won integer DEFAULT 0 NOT NULL,
    matches_scheduled integer DEFAULT 0 NOT NULL,
    seasons_played integer DEFAULT 0 NOT NULL,
    winning_ratio double precision DEFAULT 0.0 NOT NULL,
    activity_ratio double precision DEFAULT 0.0 NOT NULL,
    ranking integer,
    elo integer,
    account_id uuid NOT NULL,
    current_league_id uuid,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: refresh_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_token (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    account_id uuid NOT NULL,
    token_hash text NOT NULL,
    device_id text,
    ip_address character varying(250),
    user_agent text,
    issued_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    last_used_at timestamp with time zone,
    is_revoked boolean DEFAULT false NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: season; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.season (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    title character varying(250) NOT NULL,
    description character varying(500),
    start_date date NOT NULL,
    end_date date NOT NULL,
    creator_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: standing; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.standing (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    points integer DEFAULT 0 NOT NULL,
    matches_played integer DEFAULT 0 NOT NULL,
    matches_won integer DEFAULT 0 NOT NULL,
    sets_won integer DEFAULT 0 NOT NULL,
    sets_lost integer DEFAULT 0 NOT NULL,
    games_won integer DEFAULT 0 NOT NULL,
    games_lost integer DEFAULT 0 NOT NULL,
    season_id uuid NOT NULL,
    league_id uuid NOT NULL,
    player_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: account account_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_email_key UNIQUE (email);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);


--
-- Name: court court_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.court
    ADD CONSTRAINT court_name_key UNIQUE (name);


--
-- Name: court court_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.court
    ADD CONSTRAINT court_pkey PRIMARY KEY (id);


--
-- Name: league league_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.league
    ADD CONSTRAINT league_pkey PRIMARY KEY (id);


--
-- Name: match match_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_pkey PRIMARY KEY (id);


--
-- Name: player player_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player
    ADD CONSTRAINT player_pkey PRIMARY KEY (id);


--
-- Name: refresh_token refresh_token_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT refresh_token_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: season season_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.season
    ADD CONSTRAINT season_pkey PRIMARY KEY (id);


--
-- Name: standing standing_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.standing
    ADD CONSTRAINT standing_pkey PRIMARY KEY (id);


--
-- Name: court court_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.court
    ADD CONSTRAINT court_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.account(id);


--
-- Name: league league_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.league
    ADD CONSTRAINT league_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.account(id);


--
-- Name: league league_season_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.league
    ADD CONSTRAINT league_season_id_fkey FOREIGN KEY (season_id) REFERENCES public.season(id) ON DELETE CASCADE;


--
-- Name: match match_court_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_court_id_fkey FOREIGN KEY (court_id) REFERENCES public.court(id);


--
-- Name: match match_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.player(id) ON DELETE CASCADE;


--
-- Name: match match_league_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_league_id_fkey FOREIGN KEY (league_id) REFERENCES public.league(id) ON DELETE CASCADE;


--
-- Name: match match_player_one_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_player_one_id_fkey FOREIGN KEY (player_one_id) REFERENCES public.player(id);


--
-- Name: match match_player_two_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_player_two_id_fkey FOREIGN KEY (player_two_id) REFERENCES public.player(id);


--
-- Name: match match_season_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_season_id_fkey FOREIGN KEY (season_id) REFERENCES public.season(id) ON DELETE CASCADE;


--
-- Name: match match_winner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.match
    ADD CONSTRAINT match_winner_id_fkey FOREIGN KEY (winner_id) REFERENCES public.player(id);


--
-- Name: player player_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player
    ADD CONSTRAINT player_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE;


--
-- Name: player player_current_league_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player
    ADD CONSTRAINT player_current_league_id_fkey FOREIGN KEY (current_league_id) REFERENCES public.league(id);


--
-- Name: refresh_token refresh_token_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_token
    ADD CONSTRAINT refresh_token_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.account(id) ON DELETE CASCADE;


--
-- Name: season season_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.season
    ADD CONSTRAINT season_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.account(id);


--
-- Name: standing standing_league_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.standing
    ADD CONSTRAINT standing_league_id_fkey FOREIGN KEY (league_id) REFERENCES public.league(id) ON DELETE CASCADE;


--
-- Name: standing standing_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.standing
    ADD CONSTRAINT standing_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.player(id) ON DELETE CASCADE;


--
-- Name: standing standing_season_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.standing
    ADD CONSTRAINT standing_season_id_fkey FOREIGN KEY (season_id) REFERENCES public.season(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250106191753'),
    ('20250109114332'),
    ('20250110182359'),
    ('20250110182416'),
    ('20250110182427'),
    ('20250110191222'),
    ('20250110195840'),
    ('20250111202415');
