package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"note_app_server/global"
	"note_app_server/model/noteModel"
)

// GetNoteWithNid 获取笔记详情
func GetNoteWithNid(ctx context.Context, nid string) (*noteModel.NoteDetail, error) {
	note := new(noteModel.NoteDetail)
	sql := `select 
    			n.nid as nid, 
    			u.uid as uid, 
    			u.avatarUrl as avatarUrl, 
    			u.username as username, 
    			n.pics as pics, 
    			n.title as title, 
    			n.content as content, 
    			n.created_at as created_at, 
    			n.updated_at as updated_at, 
    			n.public as public, 
    			n.category_id as categoryId, 
    			n.tags as tags, 
    			ni.likes_count as likes_count, 
    			ni.comments_count as comments_count, 
    			ni.collections_count as collections_count, 
    			ni.shares_count as shares_count, 
    			ni.views_count as views_count 
			from notes n 
			left join user_info u on n.uid = u.uid 
			left join notes_info ni on ni.nid = n.nid 
		  	where n.status = 1 and n.nid = ?`
	if err := global.Db.WithContext(ctx).Raw(sql, nid).Scan(&note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

// UpdateNoteWithUid 更新笔记
func UpdateNoteWithUid(ctx context.Context, n *noteModel.Note) error {
	result := global.Db.WithContext(ctx).Where("nid = ? and uid = ?", n.Nid, n.Uid).Updates(&n)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("update failed")
	}
	return nil
}

// GetNoteList 查询列表
func GetNoteList(ctx context.Context, start, limit int) ([]noteModel.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []noteModel.SurfaceNote
	sql := `select 
    			n.nid as nid,
    			n.uid as uid,
    			u.username as username,
    			u.avatarUrl as avatarUrl,
    			n.cover as cover,
    			n.cover_height as cover_height,
    			n.title as title, 
    			n.public as public,
    			n.category_id as category_id,
    			n.tags as tags, 
    			ni.likes_count as likes_count
			from notes n 
			join user_info u on n.uid = u.uid 
			join notes_info ni on n.nid = ni.nid 
		 	where n.status = 1 limit ?, ?`
	res := global.Db.WithContext(ctx).Model(&noteModel.SurfaceNote{}).Raw(sql, offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("没有数据了哦")
	}
	return result, nil
}

// GetNoteListWithUid 查询用户发布的帖子
func GetNoteListWithUid(ctx context.Context, uid int64, start, limit int) ([]noteModel.SurfaceNote, error) {
	offset := (start - 1) * limit
	var result []noteModel.SurfaceNote
	sql := `select n.nid as nid,
		       n.uid as uid,
		       u.username as username,
		       u.avatarUrl as avatarUrl, 
		       n.cover as cover,
		       n.cover_height as cover_height,
		       n.title as title,
		       n.public as public,
			   n.category_id as category_id,
			   n.tags as tags,
			   ni.likes_count as like_count 
		from user_info u 
		left join notes n on n.uid = u.uid 
		left join notes_info ni on ni.nid = n.nid 
	  	where u.uid = ? and n.status = 1 limit ?, ?`
	res := global.Db.WithContext(ctx).Model(&noteModel.SurfaceNote{}).Raw(sql, uid, offset, limit).Scan(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(result) == 0 {
		return nil, errors.New("没有更多数据了哦")
	}
	return result, nil
}

// GetNoteListWithKeyword 带关键词搜索帖子
func GetNoteListWithKeyword(ctx context.Context, index, keyword string, offset, limit *int) ([]noteModel.ESNote, error) {
	analyzer := "ik_smart"
	re, err := global.ESClient.Search().Index(index).Request(&search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						MatchPhrase: map[string]types.MatchPhraseQuery{
							"all": {
								Analyzer: &analyzer,
								Query:    keyword,
							},
						},
					},
					{
						Term: map[string]types.TermQuery{
							"public": {Value: true},
						},
					},
					{
						Term: map[string]types.TermQuery{
							"status": {Value: 1},
						},
					},
				},
			},
		},
		Source_: &types.SourceFilter{
			Excludes: []string{"public", "status"},
		},
		From: offset,
		Size: limit,
	}).Do(ctx)

	if err != nil {
		return nil, err
	}
	result := make([]noteModel.ESNote, 0)
	for _, i := range re.Hits.Hits {
		var note noteModel.ESNote
		if err = json.Unmarshal(i.Source_, &note); err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, note)
	}

	return result, nil
}
