CREATE TABLE wdj_app
(
  apk          TEXT PRIMARY KEY,
  name         TEXT,
  page_url     TEXT,
  icon_url     TEXT,
  down_url     TEXT,
  size         BIGINT,
  install_cnt  BIGINT,
  comment_cnt  BIGINT,
  favor_rate   BIGINT,
  last_mtime   DATE,
  last_change  TEXT,
  last_version TEXT,
  vendor       TEXT,
  subtitle     TEXT,
  review       TEXT,
  description  TEXT,
  system       TEXT,
  permissions  TEXT [],
  tags         TEXT [],
  categories   TEXT [],
  crumb        TEXT [],
  relate_apps  TEXT [],
  screenshots  TEXT [],
  comments     TEXT [],
  crawl_time   DATE DEFAULT CURRENT_DATE
);


COMMENT ON TABLE wdj_app IS '豌豆荚应用数据表';
COMMENT ON COLUMN wdj_app.apk IS 'APK名称，即PackageName';
COMMENT ON COLUMN wdj_app.name IS '应用名称';
COMMENT ON COLUMN wdj_app.page_url IS '应用页面URL';
COMMENT ON COLUMN wdj_app.icon_url IS '应用图标(图片URL)';
COMMENT ON COLUMN wdj_app.down_url IS 'APK下载地址';
COMMENT ON COLUMN wdj_app.size IS 'APK尺寸';
COMMENT ON COLUMN wdj_app.install_cnt IS '安装数';
COMMENT ON COLUMN wdj_app.comment_cnt IS '评论数';
COMMENT ON COLUMN wdj_app.favor_rate IS '好评率';
COMMENT ON COLUMN wdj_app.last_mtime IS '最近更新时间';
COMMENT ON COLUMN wdj_app.last_change IS '最近更新内容';
COMMENT ON COLUMN wdj_app.last_version IS '最新版本';
COMMENT ON COLUMN wdj_app.vendor IS '厂商';
COMMENT ON COLUMN wdj_app.subtitle IS '副标题';
COMMENT ON COLUMN wdj_app.review IS '编辑评论';
COMMENT ON COLUMN wdj_app.description IS '应用描述';
COMMENT ON COLUMN wdj_app.system IS '系统要求';
COMMENT ON COLUMN wdj_app.permissions IS '权限要求[数组]';
COMMENT ON COLUMN wdj_app.tags IS '应用标签[数组]';
COMMENT ON COLUMN wdj_app.categories IS '应用类目[数组]';
COMMENT ON COLUMN wdj_app.crumb IS '导航栏面包屑[数组]';
COMMENT ON COLUMN wdj_app.relate_apps IS '相关应用(PkgName)[数组]';
COMMENT ON COLUMN wdj_app.screenshots IS '应用截图(图片URL)[数组]';
COMMENT ON COLUMN wdj_app.comments IS '评论(YYYYmmDD,Username,Content)[数组]';
COMMENT ON COLUMN wdj_app.crawl_time IS '爬取时间';


CREATE TABLE android (
  apk  TEXT PRIMARY KEY,
  key  TEXT,
  name TEXT,
  wdj  SMALLINT DEFAULT 0 NOT NULL,
  sjqq SMALLINT DEFAULT 0 NOT NULL,
  mi   SMALLINT DEFAULT 0 NOT NULL
);

COMMENT ON TABLE android IS '安卓应用表';
COMMENT ON COLUMN android.apk IS 'APK名称，即PackageName';
COMMENT ON COLUMN android.key IS '友盟分配的AppKey';
COMMENT ON COLUMN android.name IS '应用名称';
COMMENT ON COLUMN android.wdj IS '豌豆荚爬取状态{0:未爬,1:爬取失败,2:爬取成功}';
COMMENT ON COLUMN android.sjqq IS '应用宝爬取状态{0:未爬,1:爬取失败,2:爬取成功}';
COMMENT ON COLUMN android.mi IS '小米应用商店爬取状态{0:未爬,1:爬取失败,2:爬取成功}';


INSERT INTO wdj_app2 SELECT
                       apk,
                       name,
                       'http://www.wandoujia.com/apps/' || apk AS page_url,
                       icon                                    AS icon_url,
                       url                                     AS down_url,
                       size,
                       install_cnt,
                       comment_cnt,
                       favor_rate,
                       last_mtime,
                       last_change,
                       last_version,
                       vendor,
                       subtitle,
                       review,
                       description,
                       system,
                       permissions,
                       tags,
                       categories,
                       crumb,
                       relate_apps,
                       screenshots,
                       comments,
                       crawl_time
                     FROM wdj_app
                     ORDER BY apk ASC;