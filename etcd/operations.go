package etcd

import (
	"context"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
)

func FetchMember(
	ctx context.Context,
	c *clientv3.Client,
) (
	*model.Cluster,
	error,
) {
	var clusterState model.Cluster
	memberListResponse, err := c.MemberList(ctx)
	if err != nil {
		return nil, err
	}

	clusterState.ID = memberListResponse.Header.ClusterId
	for _, member := range memberListResponse.Members {
		statusResponse, err := c.Status(ctx, member.ClientURLs[0])
		if err != nil {
			return nil, err
		}
		clusterState.Leader = statusResponse.Leader
		clusterState.Members = append(clusterState.Members, model.ClusterMember{
			ID:         member.ID,
			Name:       member.Name,
			PeerURLS:   member.PeerURLs,
			ClientURLS: member.ClientURLs,
			IsLearner:  member.IsLearner,
			Status: model.MemberStatus{
				Version:          statusResponse.Version,
				DbSize:           statusResponse.DbSize,
				DbSizeInUse:      statusResponse.DbSizeInUse,
				RaftIndex:        statusResponse.RaftIndex,
				RaftTerm:         statusResponse.RaftTerm,
				RaftAppliedIndex: statusResponse.RaftAppliedIndex,
			},
		})
	}
	return &clusterState, err
}

func PutKV(
	ctx context.Context,
	c *clientv3.Client,
	putKVRequest *model.PutKVRequest,
) (
	*model.PutKVResponse,
	error,
) {

	putResponse, err := c.Put(ctx, putKVRequest.Key, putKVRequest.Value, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}
	kvResponse := model.PutKVResponse{
		MemberID:  putResponse.Header.GetMemberId(),
		RaftTerm:  putResponse.Header.GetRaftTerm(),
		Revision:  putResponse.Header.GetRevision(),
		ClusterID: putResponse.Header.GetClusterId(),
		NewKV:     true,
	}
	if putResponse.PrevKv != nil {
		kvResponse.NewKV = false
		kvResponse.PreviousKV = model.KV{
			Key:     string(putResponse.PrevKv.Key),
			Version: putResponse.PrevKv.Version,
			Value:   string(putResponse.PrevKv.Value),
			Lease:   putResponse.PrevKv.Lease,
		}
	}
	return &kvResponse, nil
}

func GetKV(
	ctx context.Context,
	c *clientv3.Client,
	getKVRequest *model.GetKVRequest,
) (
	*model.GetKVResponse,
	error,
) {
	var ops []clientv3.OpOption

	if getKVRequest.FromKey {
		ops = append(ops, clientv3.WithFromKey())
	}
	if getKVRequest.Limit != 0 {
		ops = append(ops, clientv3.WithLimit(getKVRequest.Limit))
	}
	if getKVRequest.KeysOnly {
		ops = append(ops, clientv3.WithKeysOnly())
	}
	if getKVRequest.Prefix {
		ops = append(ops, clientv3.WithPrefix())
	}
	if getKVRequest.Revision != 0 {
		ops = append(ops, clientv3.WithRev(getKVRequest.Revision))
	}
	if getKVRequest.Range != "" {
		ops = append(ops, clientv3.WithRange(getKVRequest.Range))
	}
	if getKVRequest.CountOnly {
		ops = append(ops, clientv3.WithCountOnly())
	}

	getResponse, err := c.Get(ctx, getKVRequest.Key, ops...)
	if err != nil {
		return nil, err
	}
	var getKVResponse = model.GetKVResponse{
		ClusterID: getResponse.Header.GetClusterId(),
		MemberID:  getResponse.Header.GetMemberId(),
		Revision:  getResponse.Header.GetRevision(),
		RaftTerm:  getResponse.Header.GetRaftTerm(),
		More:      getResponse.More,
		Count:     getResponse.Count,
	}
	for _, kv := range getResponse.Kvs {
		getKVResponse.KeyValues = append(getKVResponse.KeyValues, model.KV{
			Key:            string(kv.Key),
			Value:          string(kv.Value),
			Version:        kv.Version,
			CreateRevision: kv.CreateRevision,
			ModRevision:    kv.ModRevision,
			Lease:          kv.Lease,
		})
	}
	return &getKVResponse, nil
}

func DeleteKV(
	ctx context.Context,
	c *clientv3.Client,
	deleteKVRequest *model.DeleteKVRequest,
) (
	*model.DeleteKVResponse,
	error,
) {
	var ops []clientv3.OpOption

	if deleteKVRequest.FromKey {
		ops = append(ops, clientv3.WithFromKey())
	}
	if deleteKVRequest.Prefix {
		ops = append(ops, clientv3.WithPrefix())
	}
	if deleteKVRequest.Range != "" {
		ops = append(ops, clientv3.WithRange(deleteKVRequest.Range))
	}

	deleteResponse, err := c.Delete(ctx, deleteKVRequest.Key, ops...)
	if err != nil {
		return nil, err
	}

	return &model.DeleteKVResponse{
		ClusterID: deleteResponse.Header.GetClusterId(),
		MemberID:  deleteResponse.Header.GetMemberId(),
		Revision:  deleteResponse.Header.GetRevision(),
		RaftTerm:  deleteResponse.Header.GetRaftTerm(),
		Count:     deleteResponse.Deleted ,
	}, nil
}