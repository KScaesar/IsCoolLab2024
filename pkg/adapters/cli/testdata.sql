INSERT INTO users (username) VALUES ('user1');
INSERT INTO file_systems (id, username) VALUES ('01HYXCC8AJ35Q5KKVACBGYDF5T', 'user1');
INSERT INTO folders (id, parent_id, fs_id, name, description, created_time) VALUES ('01HYXCC8AJ35Q5KKVACDEC38G7', '', '01HYXCC8AJ35Q5KKVACBGYDF5T', '/', '', '2024-05-27 23:00:00+08:00');
INSERT INTO folders (id, parent_id, fs_id, name, description, created_time) VALUES ('01HYXCD1CD3VFFRYB9BWV19TM8', '01HYXCC8AJ35Q5KKVACDEC38G7', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder1', '', '2024-05-27 23:00:03+08:00');
INSERT INTO folders (id, parent_id, fs_id, name, description, created_time) VALUES ('01HYXCD1CGB36V08CNRGJQMZHT', '01HYXCC8AJ35Q5KKVACDEC38G7', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder2', 'qa-folder', '2024-05-27 23:00:01+08:00');
INSERT INTO folders (id, parent_id, fs_id, name, description, created_time) VALUES ('01HYXD0GV43XKBZ7Y1YDK7QDBQ', '01HYXCC8AJ35Q5KKVACDEC38G7', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder3', '', '2024-05-27 23:00:02+08:00');
INSERT INTO files (id, name, folder_id, fs_id, foldername, description, created_time) VALUES ('01HYYMFNZSFQ2FWPN1DYFTPADH', 'file1', '01HYXCD1CD3VFFRYB9BWV19TM8', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder1', '', '2024-05-27 23:00:03+08:00');
INSERT INTO files (id, name, folder_id, fs_id, foldername, description, created_time) VALUES ('01HYYMMTX8F4D2BESDCAD2YXS5', 'file2', '01HYXCD1CD3VFFRYB9BWV19TM8', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder1', 'qa-file', '2024-05-27 23:00:01+08:00');
INSERT INTO files (id, name, folder_id, fs_id, foldername, description, created_time) VALUES ('01HYYMN2H854NWJJ32HJRCQKC0', 'file3', '01HYXCD1CD3VFFRYB9BWV19TM8', '01HYXCC8AJ35Q5KKVACBGYDF5T', 'folder1', '', '2024-05-27 23:00:02+08:00');


INSERT INTO users (username) VALUES ('user2');
INSERT INTO file_systems (id, username) VALUES ('01HYXD38S85V0H1JF9CMWYBMBW', 'user2');
INSERT INTO folders (id, parent_id, fs_id, name, description, created_time) VALUES ('01HYXD4H3PPAWTFSEVTVSBKPMK', '', '01HYXD38S85V0H1JF9CMWYBMBW', '/', '', '2024-05-27 23:00:00+08:00');
