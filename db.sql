-- ************************************************************
-- Sequel Pro SQL dump
-- Version 4541
--
-- http://www.sequelpro.com/
-- https://github.com/sequelpro/sequelpro
--
-- Host: 127.0.0.1 (MySQL 5.7.29)
-- Database: data-sync
-- Generation Time: 2021-06-06 05:08:10 +0000
-- ************************************************************


/* SQLINES DEMO *** ARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/* SQLINES DEMO *** ARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/* SQLINES DEMO *** LLATION_CONNECTION=@@COLLATION_CONNECTION */;
/* SQLINES DEMO *** tf8 */;
/* SQLINES DEMO *** REIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/* SQLINES DEMO *** L_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/* SQLINES DEMO *** L_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dump of table beacon_blocks
-- ------------------------------------------------------------
  
-- ------------------------------------------------------------

DROP TABLE IF EXISTS p_tokens;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
--CREATE SEQUENCE p_tokens_seq;

CREATE TABLE p_tokens (
  id  SERIAL PRIMARY KEY, 
  token_id varchar(255) DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  symbol varchar(255) DEFAULT NULL,
  decimal double precision DEFAULT NULL,
  price double precision DEFAULT NULL 
) ;



-- Dump of table pde_pool_pairs
-- ------------------------------------------------------------

DROP TABLE IF EXISTS pde_pool_pairs;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE pde_pool_pairs (
  id  SERIAL PRIMARY KEY, 
  token1_id_str varchar(255) DEFAULT NULL,
  token1_pool_value bigint DEFAULT NULL,
  token2_id_str varchar(255) DEFAULT NULL,
  token2_pool_value bigint   DEFAULT NULL,
  token1_to_token2_price bigint  DEFAULT NULL,
  token2_to_token1_price bigint  DEFAULT NULL,
  beacon_height bigint  DEFAULT NULL,
  beacon_time_stamp timestamp(0) DEFAULT NULL 
) ;



-- Dump of table pde_trades
-- ------------------------------------------------------------

 DROP TABLE IF EXISTS pde_trades;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE pde_trades ( 
  id  SERIAL PRIMARY KEY, 
  trader_address_str varchar(255) DEFAULT NULL,
  receiving_tokenid_str varchar(255) DEFAULT NULL,
  receive_amount bigint check (receive_amount > 0) DEFAULT NULL,
  token1_id_str varchar(255) DEFAULT NULL,
  token2_id_str varchar(255) DEFAULT NULL,
  shard_id smallint check (shard_id > 0) DEFAULT NULL,
  requested_tx_id varchar(255) DEFAULT NULL,
  status varchar(255) DEFAULT NULL,
  beacon_height bigint check (beacon_height > 0) DEFAULT NULL,
  beacon_time_stamp timestamp(0) DEFAULT NULL,
  price double precision DEFAULT NULL,
  amount double precision DEFAULT NULL
) ;

  
-- Dump of table short_transactions
-- ------------------------------------------------------------

DROP TABLE IF EXISTS short_transactions;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE short_transactions (
   id  SERIAL PRIMARY KEY, 
  tx_id varchar(255) DEFAULT NULL,
  data text,
  proof text[],
  proof_detail text,
  metadata text,
  transacted_privacy_coin text,
  transacted_privacy_coin_proof_detail text,
  meta_data_type int DEFAULT NULL,
  created_time timestamp(0) DEFAULT NULL
) ;



-- Dump of table tokens
-- ------------------------------------------------------------

DROP TABLE IF EXISTS tokens;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
--CREATE SEQUENCE tokens_seq;

CREATE TABLE tokens (
  id  SERIAL PRIMARY KEY, 
  token_id varchar(255) DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  symbol varchar(255) DEFAULT NULL,
  count_tx int DEFAULT NULL,
  supply bigint check (supply > 0) DEFAULT NULL,
  list_hash_tx json DEFAULT NULL,
  data varchar(255) DEFAULT NULL,
  info varchar(255) DEFAULT NULL 
) ;



-- Dump of table transactions
-- ------------------------------------------------------------
DROP TABLE IF EXISTS transactions;

-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE transactions (
  id  SERIAL PRIMARY KEY, 
  tx_id varchar(255) DEFAULT NULL,
  tx_version int DEFAULT NULL,
  tx_type varchar(255) DEFAULT NULL,
  data text DEFAULT NULL,
  shard_id int DEFAULT NULL,
  prv_fee bigint DEFAULT NULL,
  info text DEFAULT NULL,
  proof text,
  proof_detail text,
  metadata text,
  transacted_privacy_coin text,
  transacted_privacy_coin_proof_detail text,
  transacted_privacy_coin_fee bigint DEFAULT NULL,
  created_time timestamp(0) DEFAULT NULL,
  block_height bigint DEFAULT NULL,
  block_hash varchar(255) DEFAULT NULL,
  meta_data_type int DEFAULT NULL,
  public_key_list varchar[],
  serial_number_list varchar[],
  coin_commitment_list varchar[],
  shield_type int DEFAULT NULL, 
  amount_shield double precision DEFAULT NULL,
  price double precision DEFAULT NULL,
  token_id varchar(255) DEFAULT NULL,
  token_name varchar(255) DEFAULT NULL 
) ; 

DROP TABLE public.shard_blocks;

CREATE TABLE public.shard_blocks
(
    id  SERIAL PRIMARY KEY, 
    block_hash character varying(255) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    block_version integer,
    block_height bigint,
    block_producer character varying(255) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    epoch bigint,
    round integer,
    created_time timestamp(0) without time zone DEFAULT NULL::timestamp without time zone,
    data text COLLATE pg_catalog."default",
    count_tx integer,
    shard_id integer,
    pre_block character varying(255) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    next_block character varying(255) COLLATE pg_catalog."default" DEFAULT NULL::character varying,
    list_hash_tx json,
    beacon_block_height bigint 
)
