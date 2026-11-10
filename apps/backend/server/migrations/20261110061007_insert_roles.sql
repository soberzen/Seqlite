-- +goose Up
INSERT INTO roles (code, name, description) VALUES
('owner',     '所有者',   '空间/文档的创建者，拥有全部权限'),
('admin',     '管理员',   '管理成员和内容，但不能删除空间/文档'),
('editor',    '编辑者',   '可查看和编辑内容，不能管理成员'),
('viewer',    '查看者',   '只能查看，不能修改'),
('commenter', '评论者',   '只能查看和评论')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DELETE FROM roles WHERE code IN ('owner','admin','editor','viewer','commenter');
