DROP TRIGGER IF EXISTS trg_create_tenant_ledger_accounts ON tenants;
DROP FUNCTION IF EXISTS create_tenant_ledger_accounts();
DROP TRIGGER IF EXISTS trg_ledger_balance ON ledger_entries;
DROP FUNCTION IF EXISTS verify_ledger_balance();
DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS ledger_transactions;
DROP TABLE IF EXISTS ledger_accounts;
DROP TYPE IF EXISTS ledger_account_scope;
DROP TYPE IF EXISTS ledger_account_kind;