DELETE FROM "gateways" WHERE id=1;
DELETE FROM "gateways" WHERE id=2;
DELETE FROM "gateways" WHERE id=3;


DELETE FROM "config" WHERE key='GrpcVersion';
DELETE FROM "config" WHERE key='RestVersion';
DELETE FROM "config" WHERE key='DispatcherVersion';

drop sequence Seq_Account;