INSERT INTO "gateways"(
    id, user_name, hash_password, ips, gateway_name, status)
VALUES (1, 'company1', '927258103250a45487396663c039438732391ab8e843961ae0470d00ecb52379', 'localhost, 127.0.0.1', 'WEB', true);

INSERT INTO "gateways"(
    id, user_name, hash_password, ips, gateway_name, status)
VALUES (2, 'company2', '927258103250a45487396663c039438732391ab8e843961ae0470d00ecb52379', 'localhost, 127.0.0.1', 'Android', true);

INSERT INTO "gateways"(
    id, user_name, hash_password, ips, gateway_name, status)
VALUES (3, 'company3', '927258103250a45487396663c039438732391ab8e843961ae0470d00ecb52379', 'localhost, 127.0.0.1', 'IOS', true);

INSERT INTO "config"(key, value)VALUES ('GrpcVersion', '1.0.0');
INSERT INTO "config"(key, value)VALUES ('RestVersion', '1.0.0');
INSERT INTO "config"(key, value)VALUES ('DispatcherVersion', '1.0.0');

CREATE SEQUENCE Seq_Account
    INCREMENT 5
    START 100;