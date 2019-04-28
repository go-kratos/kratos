package service

//func (s *Service) updateTelDuplicateJob() {
//	log.Info("update tel duplicate job start...")
//	res, err := s.d.UserTelDuplicate(context.Background())
//	if err != nil {
//		log.Error("fail to get UserTelDuplicate, error(%+v)", err)
//		return
//	}
//	for _, r := range res {
//		userTel := &model.UserTel{
//			Mid: r.Mid,
//			Tel: r.Tel,
//			Cid: r.Cid,
//		}
//		if _, err = s.d.UpdateUserTel(context.Background(), userTel); err != nil {
//			log.Error("fail to update user tel userTel(%+v), error(%+v)", userTel, err)
//			continue
//		}
//		log.Info("update userTel by userTelDuplicate success, userTelDuplicate(%+v)", r)
//		if _, err = s.d.UpdateUserTelDuplicateStatus(context.Background(), r.Mid); err != nil {
//			log.Error("fail to update user tel duplicate status, mid(%d) error(%+v)", r.Mid, err)
//		}
//		log.Info("update userTelDuplicate status success, mid(%d)", r.Mid)
//	}
//	log.Info("update tel duplicate job end...")
//}
//
//func (s *Service) updateEmailDuplicateJob() {
//	log.Info("update email duplicate job start...")
//	res, err := s.d.UserEmailDuplicate(context.Background())
//	if err != nil {
//		log.Error("fail to get UserTelDuplicate, error(%+v)", err)
//		return
//	}
//	for _, r := range res {
//		userEmail := &model.UserEmail{
//			Mid:   r.Mid,
//			Email: r.Email,
//		}
//		if _, err = s.d.UpdateUserEmail(context.Background(), userEmail); err != nil {
//			log.Error("fail to update user email userEmail(%+v), error(%+v)", userEmail, err)
//			continue
//		}
//		log.Info("update userEmail by userEmailDuplicate success, userEmailDuplicate(%+v)", r)
//		if _, err = s.d.UpdateUserTelDuplicateStatus(context.Background(), r.Mid); err != nil {
//			log.Error("fail to update user email duplicate status, mid(%d) error(%+v)", r.Mid, err)
//		}
//		log.Info("update userEmailDuplicate status success, mid(%d)", r.Mid)
//	}
//	log.Info("update email duplicate job end...")
//}

/**
func (s *Service) checkTelDuplicateJob() {
	log.Info("check tel duplicate job start...")
	res, err := s.d.UserTelDuplicate(context.Background())
	if err != nil {
		log.Error("fail to get UserTelDuplicate, error(%+v)", err)
		return
	}
	var (
		asoAccount *model.OriginAccount
		userTel    *model.UserTel
	)
	for _, r := range res {
		if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check tel duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		if userTel, err = s.d.GetUserTelByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check tel duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		a := s.convertAccountToUserTel(asoAccount)
		originHex := hex.EncodeToString(a.Tel)
		newHex := hex.EncodeToString(userTel.Tel)
		if originHex == newHex {
			log.Info("check user tel duplicate success, userTelDuplicate(%+v)", r)
			if _, err = s.d.UpdateUserTelDuplicateStatus(context.Background(), r.ID); err != nil {
				log.Error("fail to update user tel duplicate status, id(%d) error(%+v)", r.ID, err)
			}
		} else {
			log.Info("fail to check user tel duplicate, new(%s) origin(%s) userTelDuplicate(%+v)", newHex, originHex, r)
		}
	}
	log.Info("update tel duplicate job end...")
}

func (s *Service) checkEmailDuplicateJob() {
	log.Info("check email duplicate job start...")
	res, err := s.d.UserEmailDuplicate(context.Background())
	if err != nil {
		log.Error("fail to get UserEmailDuplicate, error(%+v)", err)
		return
	}
	var (
		asoAccount *model.OriginAccount
		userEmail  *model.UserEmail
	)
	for _, r := range res {
		if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check email duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		if userEmail, err = s.d.GetUserEmailByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check email duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		a := s.convertAccountToUserEmail(asoAccount)
		originHex := hex.EncodeToString(a.Email)
		newHex := hex.EncodeToString(userEmail.Email)
		if originHex == newHex {
			log.Info("check user email duplicate success, userEmailDuplicate(%+v)", r)
			if _, err = s.d.UpdateUserTelDuplicateStatus(context.Background(), r.ID); err != nil {
				log.Error("fail to update user email duplicate status, id(%d) error(%+v)", r.ID, err)
			}
		} else {
			log.Info("fail to check user tel duplicate, new(%s) origin(%s) userEmailDuplicate(%+v)", newHex, originHex, r)
		}
	}
	log.Info("update email duplicate job end...")
}**/
